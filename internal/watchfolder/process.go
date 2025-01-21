package watchfolder

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/database/repository"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/internal/service"
	"github.com/welovemedia/ffmate/sev"
	"github.com/yosev/debugo"
)

type Watchfolder struct {
	Sev                   *sev.Sev
	WatchfolderService    *service.WatchfolderService
	WatchfolderRepository *repository.Watchfolder
	WebhookService        *service.WebhookService
	WebsocketService      *service.WebsocketService
	TaskService           *service.TaskService
}

var debug = debugo.New("watchfolder")

type fileState struct {
	Size     int64
	Attempts int
}

func (w *Watchfolder) Init() {
	watchfolders, total, _ := w.WatchfolderRepository.List(-1, -1)
	debug.Debugf("initializing %d watchfolders", total)

	for _, watchfolder := range *watchfolders {
		go w.process(&watchfolder)
	}
}

func (w *Watchfolder) process(watchfolder *model.Watchfolder) {
	fileStates := make(map[string]*fileState)
	processedFiles := make(map[string]bool)
	var mu sync.Mutex

	for {
		// Walk the directory
		err := filepath.Walk(watchfolder.Path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Skip directories
			if info.IsDir() {
				return nil
			}

			// Skip invisible files
			if strings.HasPrefix(filepath.Base(path), ".") {
				return nil
			}

			mu.Lock()
			defer mu.Unlock()

			// Check if the file has already been processed
			if _, seen := processedFiles[path]; seen {
				return nil
			}

			// Determine if the file is ready for processing
			if shouldProcessFile(path, info, fileStates, watchfolder.GrowthChecks) {
				w.createTask(path, watchfolder)
				processedFiles[path] = true
				delete(fileStates, path) // Remove from tracking
			}

			return nil
		})

		if err != nil {
			w.Sev.Logger().Errorf("error walking watchfolder directory (uuid: %s): %v", watchfolder.Uuid, err)
		}

		time.Sleep(time.Duration(watchfolder.Interval))
	}
}

func (w *Watchfolder) createTask(path string, watchfolder *model.Watchfolder) {
	debug.Debugf("created new task for file: %s", path)
	w.TaskService.NewTask(&dto.NewTask{
		Preset:    watchfolder.Preset,
		Name:      filepath.Base(path),
		InputFile: path,
	}, "")
}

// shouldProcessFile determines if a file is ready for processing based on growth attempts.
func shouldProcessFile(path string, info os.FileInfo, fileStates map[string]*fileState, growthChecks int) bool {
	if growthChecks == 0 {
		// If no growth checks are required, the file is ready immediately
		return true
	}

	// Get or initialize the file state
	state, exists := fileStates[path]
	if !exists {
		fileStates[path] = &fileState{Size: info.Size(), Attempts: 1}
		return false
	}

	// Check if the file size is stable
	if info.Size() == state.Size {
		state.Attempts++
		if state.Attempts >= growthChecks {
			return true
		}
	} else {
		// File size changed, reset attempts
		state.Size = info.Size()
		state.Attempts = 1
	}

	return false
}
