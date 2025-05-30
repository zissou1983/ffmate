package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/mattn/go-shellwords"
	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/database/repository"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/internal/ffmpeg"
	"github.com/welovemedia/ffmate/internal/service"
	"github.com/welovemedia/ffmate/internal/utils/wildcards"
	"github.com/welovemedia/ffmate/sev"
	"github.com/yosev/debugo"
)

type Queue struct {
	Sev                *sev.Sev
	TaskRepository     *repository.Task
	MaxConcurrentTasks uint
}

var debug = debugo.New("queue")

var (
	taskCtx = make(map[string]context.CancelCauseFunc)
	taskMu  = &sync.Mutex{}
)

func (q *Queue) Init() {
	go func() {
		for {
			taskMu.Lock()
			if len(taskCtx) < int(q.MaxConcurrentTasks) {
				task, err := q.TaskRepository.NextQueued()
				if err != nil {
					q.Sev.Logger().Errorf("failed to receive queued task from db: %v", err)
				} else if task == nil {
					debug.Debug("no queued tasks found")
				} else {
					ctx, cancelTask := context.WithCancelCause(context.Background())
					taskCtx[task.Uuid] = cancelTask
					go q.processTask(task, ctx, func() {
						taskMu.Lock()
						defer taskMu.Unlock()
						delete(taskCtx, task.Uuid)
					})
				}
			} else {
				debug.Debugf("maximum concurrent tasks reached (tasks: %d/%d)", len(taskCtx), q.MaxConcurrentTasks)
			}
			taskMu.Unlock()
			time.Sleep(1 * time.Second)
		}
	}()
	go func() {
		for t := range service.TaskService().GetTaskUpdates() {
			taskMu.Lock()
			if fn, ok := taskCtx[t.Uuid]; ok {
				fn(errors.New("task canceled by user"))
			} else {
				q.Sev.Logger().Warnf("task not found to cancel (uuid: %s)", t.Uuid)
			}
			taskMu.Unlock()
		}
	}()
}

func (q *Queue) processTask(task *model.Task, ctx context.Context, doneFunc func()) {
	defer doneFunc()

	task.StartedAt = time.Now().UnixMilli()
	q.Sev.Logger().Infof("processing task (uuid: %s)", task.Uuid)

	err := q.prePostProcessTask(task, task.PreProcessing, "pre")
	if err != nil {
		q.failTask(task, fmt.Errorf("PreProcessing failed: %v", err))
		return
	}

	// resolve wildcards
	inFile := wildcards.Replace(task.InputFile.Raw, task.InputFile.Raw, task.OutputFile.Raw, task.Source)
	outFile := wildcards.Replace(task.OutputFile.Raw, task.InputFile.Raw, task.OutputFile.Raw, task.Source)
	task.InputFile.Resolved = inFile
	task.OutputFile.Resolved = outFile
	task.Command.Resolved = wildcards.Replace(task.Command.Raw, inFile, outFile, task.Source)
	task.Status = dto.RUNNING
	q.updateTask(task)

	q.Sev.Logger().Infof("starting processing (uuid: %s)", task.Uuid)
	err = ffmpeg.Execute(
		&ffmpeg.ExecutionRequest{
			Task:    task,
			Command: task.Command.Resolved,
			Logger:  q.Sev.Logger(),
			Ctx:     ctx,
			UpdateFunc: func(progress float64, remaining float64) {
				task.Progress = progress
				task.Remaining = remaining
				q.updateTask(task)
			},
		},
	)

	// task is done (successful or not)
	task.Progress = 100
	task.Remaining = -1

	if err != nil {
		q.Sev.Logger().Errorf("finished processing with error (uuid: %s): %v", task.Uuid, err)
		if context.Cause(ctx) != nil {
			q.cancelTask(task, context.Cause(ctx))
			return
		}
		q.failTask(task, err)
		return
	}

	q.Sev.Logger().Infof("finished processing (uuid: %s)", task.Uuid)

	err = q.prePostProcessTask(task, task.PostProcessing, "post")
	if err != nil {
		q.failTask(task, fmt.Errorf("PostProcessing failed: %v", err))
		return
	}

	task.FinishedAt = time.Now().UnixMilli()
	task.Status = dto.DONE_SUCCESSFUL
	q.updateTask(task)
	q.Sev.Logger().Infof("task successful (uuid: %s)", task.Uuid)
}

func (q *Queue) prePostProcessTask(task *model.Task, processor *dto.PrePostProcessing, processorType string) error {
	if processor != nil && (processor.SidecarPath != nil || processor.ScriptPath != nil) {
		if processorType == "pre" {
			q.Sev.Metrics().GaugeVec("task.preProcessing").WithLabelValues(strconv.FormatBool(processor.SidecarPath != nil && processor.SidecarPath.Raw == ""), strconv.FormatBool(processor.ScriptPath != nil && processor.ScriptPath.Raw == "")).Inc()
		} else {
			q.Sev.Metrics().GaugeVec("task.postProcessing").WithLabelValues(strconv.FormatBool(processor.SidecarPath != nil && processor.SidecarPath.Raw == ""), strconv.FormatBool(processor.ScriptPath != nil && processor.ScriptPath.Raw == "")).Inc()
		}
		q.Sev.Logger().Infof("starting %sProcessing (uuid: %s)", processorType, task.Uuid)
		processor.StartedAt = time.Now().UnixMilli()
		if processorType == "pre" {
			task.Status = dto.PRE_PROCESSING
		} else {
			task.Status = dto.POST_PROCESSING
		}
		q.updateTask(task)
		if processor.SidecarPath != nil && processor.SidecarPath.Raw != "" {
			b, err := json.Marshal(task.ToDto())
			if err != nil {
				q.Sev.Logger().Errorf("failed to marshal task to write sidecar file: %v", err)
			} else {
				if processorType == "pre" {
					processor.SidecarPath.Resolved = wildcards.Replace(processor.SidecarPath.Raw, task.InputFile.Raw, task.OutputFile.Raw, task.Source)
				} else {
					processor.SidecarPath.Resolved = wildcards.Replace(processor.SidecarPath.Raw, task.InputFile.Resolved, task.OutputFile.Resolved, task.Source)
				}
				q.updateTask(task)
				err = os.WriteFile(processor.SidecarPath.Resolved, b, 0644)
				if err != nil {
					processor.Error = fmt.Errorf("failed to write sidecar: %v", err).Error()
					q.Sev.Logger().Errorf("failed to write sidecar file: %v", err)
				} else {
					debug.Debugf("wrote sidecar file (uuid: %s)", task.Uuid)
				}
			}
		}

		if processor.Error == "" && processor.ScriptPath != nil && processor.ScriptPath.Raw != "" {
			if processorType == "pre" {
				processor.ScriptPath.Resolved = wildcards.Replace(processor.ScriptPath.Raw, task.InputFile.Raw, task.OutputFile.Raw, task.Source)
			} else {
				processor.ScriptPath.Resolved = wildcards.Replace(processor.ScriptPath.Raw, task.InputFile.Resolved, task.OutputFile.Resolved, task.Source)
			}
			q.updateTask(task)
			args, err := shellwords.NewParser().Parse(processor.ScriptPath.Resolved)
			if err != nil {
				processor.Error = err.Error()
				q.Sev.Logger().Errorf("failed to parse %sProcessing script (uuid: %s): %v", processorType, task.Uuid, err)
			} else {
				cmd := exec.Command(args[0], args[1:]...)
				debug.Debugf("triggered %sProcessing script (uuid: %s)", processorType, task.Uuid)

				if err := cmd.Start(); err != nil {
					processor.Error = err.Error()
					q.Sev.Logger().Errorf("failed to start %sProcessing script (uuid: %s): %v", processorType, task.Uuid, err)
				} else {
					if err := cmd.Wait(); err != nil {
						processor.Error = err.Error()
						q.Sev.Logger().Errorf("failed %sProcessing script (uuid: %s): %v", processorType, task.Uuid, err)
					}

				}
			}
		}

		processor.FinishedAt = time.Now().UnixMilli()
		if processor.Error != "" {
			q.Sev.Logger().Infof("finished %sProcessing with error (uuid: %s)", processorType, task.Uuid)
			return errors.New(processor.Error)
		}
		q.Sev.Logger().Infof("finished %sProcessing (uuid: %s)", processorType, task.Uuid)
	}
	return nil
}

func (q *Queue) cancelTask(task *model.Task, err error) {
	task.FinishedAt = time.Now().UnixMilli()
	task.Progress = 100
	task.Status = dto.DONE_CANCELED
	task.Error = err.Error()
	q.updateTask(task)
	q.Sev.Logger().Warnf("task canceled (uuid: %s): %v", task.Uuid, err)
}

func (q *Queue) failTask(task *model.Task, err error) {
	task.FinishedAt = time.Now().UnixMilli()
	task.Progress = 100
	task.Status = dto.DONE_ERROR
	task.Error = err.Error()
	q.updateTask(task)
	q.Sev.Logger().Warnf("task failed (uuid: %s):\n%v", task.Uuid, err)
}

func (q *Queue) updateTask(task *model.Task) {
	service.TaskService().UpdateTask(task)
}
