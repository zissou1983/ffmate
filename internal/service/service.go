package service

import (
	"github.com/welovemedia/ffmate/internal/database/repository"
	"github.com/welovemedia/ffmate/sev"
)

type service struct {
	preset      *presetSvc
	task        *taskSvc
	watchfolder *watchfolderSvc
	webhook     *webhookSvc
	websocket   *websocketSvc
}

var services *service

func Init(s *sev.Sev) {
	services = &service{
		preset:      &presetSvc{sev: s, presetRepository: &repository.Preset{DB: s.DB()}},
		task:        &taskSvc{sev: s, taskRepository: &repository.Task{DB: s.DB()}},
		watchfolder: &watchfolderSvc{sev: s, watchfolderRepository: &repository.Watchfolder{DB: s.DB()}},
		webhook:     &webhookSvc{sev: s, webhookRepository: &repository.Webhook{DB: s.DB()}},
		websocket:   &websocketSvc{},
	}
}

// Accessor methods
func PresetService() *presetSvc {
	return services.preset
}

func TaskService() *taskSvc {
	return services.task
}

func WatchfolderService() *watchfolderSvc {
	return services.watchfolder
}

func WebhookService() *webhookSvc {
	return services.webhook
}

func WebsocketService() *websocketSvc {
	return services.websocket
}
