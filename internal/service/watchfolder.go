package service

import (
	"errors"

	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/database/repository"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/sev"
)

type WatchfolderService struct {
	Sev                   *sev.Sev
	WatchfolderRepository *repository.Watchfolder
	WebhookService        *WebhookService
	PresetService         *PresetService
	WebsocketService      *WebsocketService
}

var watchfolderUpdates = make(chan *model.Watchfolder, 100)

func (s *WatchfolderService) GetWatchfolderUpdates() chan *model.Watchfolder {
	return watchfolderUpdates
}

func (s *WatchfolderService) ListWatchfolders(page int, perPage int) (*[]model.Watchfolder, int64, error) {
	return s.WatchfolderRepository.List(page, perPage)
}

func (s *WatchfolderService) GetWatchfolderById(uuid string) (*model.Watchfolder, error) {
	return s.WatchfolderRepository.First(uuid)
}

func (s *WatchfolderService) DeleteWatchfolder(uuid string) error {
	w, err := s.WatchfolderRepository.First(uuid)
	if err != nil {
		return err
	}

	if w.Uuid == "" {
		return errors.New("watchfolder for given uuid not found")
	}

	err = s.WatchfolderRepository.Delete(w)
	if err != nil {
		s.Sev.Logger().Warnf("failed to delete watchfolder (uuid: %s): %+v", w.Uuid, err)
		return err
	}

	s.Sev.Logger().Infof("deleted watchfolder (uuid: %s)", w.Uuid)
	watchfolderUpdates <- w

	s.Sev.Metrics().Gauge("watchfolder.deleted").Inc()
	s.WebhookService.Fire(dto.WATCHFOLDER_DELETED, w.ToDto())
	s.WebsocketService.Broadcast(WATCHFOLDER_DELETED, w.ToDto())

	return nil
}

func (s *WatchfolderService) NewWatchfolder(newWatchfolder *dto.NewWatchfolder) (*model.Watchfolder, error) {
	_, err := s.PresetService.FindByUuid(newWatchfolder.Preset)
	if err != nil {
		return nil, err
	}
	w, err := s.WatchfolderRepository.Create(newWatchfolder)

	s.Sev.Logger().Infof("created new watchfolder (uuid: %s)", w.Uuid)
	watchfolderUpdates <- w

	s.Sev.Metrics().Gauge("watchfolder.created").Inc()
	s.WebhookService.Fire(dto.WATCHFOLDER_CREATED, w.ToDto())
	s.WebsocketService.Broadcast(WATCHFOLDER_CREATED, w.ToDto())

	return w, err
}
