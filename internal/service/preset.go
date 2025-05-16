package service

import (
	"errors"

	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/database/repository"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/sev"
)

type presetSvc struct {
	service
	sev              *sev.Sev
	presetRepository *repository.Preset
}

func (s *presetSvc) FindByUuid(uuid string) (*model.Preset, error) {
	w, err := s.presetRepository.FindByUuid(uuid)
	if err != nil {
		return nil, err
	}

	if w.Uuid == "" {
		return nil, errors.New("preset for given uuid not found")
	}

	return w, nil
}

func (s *presetSvc) ListPresets(page int, perPage int) (*[]model.Preset, int64, error) {
	return s.presetRepository.List(page, perPage)
}

func (s *presetSvc) DeletePreset(uuid string) error {
	w, err := s.presetRepository.First(uuid)
	if err != nil {
		return err
	}

	if w.Uuid == "" {
		return errors.New("preset for given uuid not found")
	}

	err = s.presetRepository.Delete(w)
	if err != nil {
		s.sev.Logger().Warnf("failed to delete preset (uuid: %s): %+v", w.Uuid, err)
		return err
	}

	s.sev.Logger().Infof("deleted preset (uuid: %s)", w.Uuid)

	s.sev.Metrics().Gauge("preset.deleted").Inc()
	WebhookService().Fire(dto.PRESET_DELETED, w.ToDto())
	WebsocketService().Broadcast(PRESET_DELETED, w.ToDto())

	return nil
}

func (s *presetSvc) NewPreset(newPreset *dto.NewPreset) (*model.Preset, error) {
	w, err := s.presetRepository.Create(newPreset)
	s.sev.Logger().Infof("created new preset (uuid: %s)", w.Uuid)

	s.sev.Metrics().Gauge("preset.created").Inc()
	WebhookService().Fire(dto.PRESET_CREATED, w.ToDto())
	WebsocketService().Broadcast(PRESET_CREATED, w.ToDto())

	return w, err
}

func (s *presetSvc) UpdatePreset(presetUuid string, newPreset *dto.NewPreset) (*model.Preset, error) {
	p, err := s.FindByUuid(presetUuid)
	if err != nil {
		return nil, err
	}

	p.Name = newPreset.Name
	p.Description = newPreset.Description
	p.Command = newPreset.Command
	p.PreProcessing = newPreset.PreProcessing
	p.PostProcessing = newPreset.PostProcessing
	p.OutputFile = newPreset.OutputFile
	p.Priority = newPreset.Priority

	err = s.presetRepository.Update(p)
	if err != nil {
		s.sev.Logger().Warnf("failed to update preset (uuid: %s): %+v", p.Uuid, err)
		return nil, err
	}

	s.sev.Metrics().Gauge("preset.updated").Inc()
	WebhookService().Fire(dto.PRESET_UPDATED, p.ToDto())
	WebsocketService().Broadcast(PRESET_UPDATED, p.ToDto())

	return p, err
}
