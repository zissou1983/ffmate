package service

import (
	"errors"

	"github.com/welovemedia/ffmate/internal/database/model"
	"github.com/welovemedia/ffmate/internal/database/repository"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/sev"
)

type webhookSvc struct {
	service
	sev               *sev.Sev
	webhookRepository *repository.Webhook
}

func (s *webhookSvc) ListWebhooks(page int, perPage int) (*[]model.Webhook, int64, error) {
	return s.webhookRepository.List(page, perPage)
}

func (s *webhookSvc) DeleteWebhook(uuid string) error {
	w, err := s.webhookRepository.First(uuid)
	if err != nil {
		return err
	}

	if w.Uuid == "" {
		return errors.New("webhook for given uuid not found")
	}

	err = s.webhookRepository.Delete(w)
	if err != nil {
		s.sev.Logger().Warnf("failed to delete webhook for event %s (uuid: %s): %+v", w.Event, w.Uuid, err)
		return err
	}

	s.sev.Logger().Infof("deleted webhook for event %s (uuid: %s)", w.Event, w.Uuid)

	s.sev.Metrics().Gauge("webhook.deleted").Inc()
	s.Fire(dto.WEBHOOK_DELETED, w)

	return nil
}

func (s *webhookSvc) NewWebhook(webhook *dto.NewWebhook) (*model.Webhook, error) {
	w, err := s.webhookRepository.Create(webhook.Event, webhook.Url)
	s.sev.Logger().Infof("created new webhook for event %s (uuid: %s)", w.Event, w.Uuid)

	s.sev.Metrics().Gauge("webhook.created").Inc()
	s.Fire(dto.WEBHOOK_CREATED, w)

	return w, err
}

func (s *webhookSvc) Fire(event dto.WebhookEvent, data interface{}) error {
	webhooks, err := s.webhookRepository.ListByEvent(event)
	for _, webhook := range *webhooks {
		go s.sev.FireWebhook(&webhook, data)
		s.sev.Metrics().Gauge("webhook.executed").Inc()
	}
	return err
}
