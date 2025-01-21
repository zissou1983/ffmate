package watchfolder

import (
	"crypto/md5"
	"encoding/hex"
	"io/fs"
	"path/filepath"
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

type monitor struct {
	check     int
	processed bool
}

var debug = debugo.New("watchfolder")

var watchfolderMonitor = make(map[string]*monitor)

func (w *Watchfolder) Init() {
	watchfolders, total, _ := w.WatchfolderRepository.List(-1, -1)
	debug.Debugf("initializing %d watchfolders", total)

	for _, watchfolder := range *watchfolders {
		go func() {
			for {
				w.process(&watchfolder)
				time.Sleep(time.Duration(watchfolder.Interval) * time.Second)
			}
		}()
	}
}

func (w *Watchfolder) process(watchfolder *model.Watchfolder) {
	debug.Debugf("processing watchfolder (uuid: %s)", watchfolder.Uuid)
	filepath.WalkDir(watchfolder.Path, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			id := getMD5Hash(watchfolder.Uuid + ":" + path)
			if watchfolder.GrowthCheck <= 0 {
				watchfolderMonitor[id] = &monitor{
					check:     0,
					processed: true,
				}
				w.TaskService.NewTask(&dto.NewTask{
					Preset: watchfolder.Preset,
					InputFile: path,
					Name: watchfolder.Name,
				}, "")
			} else {
				if watchfolderMonitor[id] = &monitor{
					check:     0,
					processed: true,
				}
			}
		}
		return nil
	})

}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
