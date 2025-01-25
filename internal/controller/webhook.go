package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/internal/interceptor"
	"github.com/welovemedia/ffmate/internal/service"
	"github.com/welovemedia/ffmate/sev"
	"github.com/welovemedia/ffmate/sev/exceptions"
)

type WebhookController struct {
	sev.Controller
	sev *sev.Sev

	Prefix string
}

func (c *WebhookController) Setup(s *sev.Sev) {
	c.sev = s
	s.Gin().DELETE(c.Prefix+c.getEndpoint()+"/:uuid", c.deleteWebhook)
	s.Gin().POST(c.Prefix+c.getEndpoint(), c.addWebhook)
	s.Gin().GET(c.Prefix+c.getEndpoint(), interceptor.PageLimit, c.listWebhooks)
}

// @Summary Delete a webhook
// @Description Delete a webhook by its uuid
// @Tags webhooks
// @Param uuid path string true "the webhooks uuid"
// @Produce json
// @Success 204
// @Router /webhooks/{uuid} [delete]
func (c *WebhookController) deleteWebhook(gin *gin.Context) {
	uuid := gin.Param("uuid")
	err := service.WebhookService().DeleteWebhook(uuid)

	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err))
		return
	}

	gin.AbortWithStatus(204)
}

// @Summary List all webhooks
// @Description List all existing webhooks
// @Tags webhooks
// @Produce json
// @Success 200 {object} []dto.Webhook
// @Router /webhooks [get]
func (c *WebhookController) listWebhooks(gin *gin.Context) {
	webhooks, total, err := service.WebhookService().ListWebhooks(gin.GetInt("page"), gin.GetInt("perPage"))
	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err))
		return
	}

	gin.Header("X-Total", fmt.Sprintf("%d", total))

	// Transform each webhook to its DTO
	var webhooksDTOs = []dto.Webhook{}
	for _, webhook := range *webhooks {
		webhooksDTOs = append(webhooksDTOs, *webhook.ToDto())
	}

	gin.JSON(200, webhooksDTOs)
}

// @Summary Add a new webhook
// @Description Add a new webhook for an event
// @Tags webhooks
// @Accept json
// @Param request body dto.NewWebhook true "new webhook"
// @Produce json
// @Success 200 {object} dto.Webhook
// @Router /webhooks [post]
func (c *WebhookController) addWebhook(gin *gin.Context) {
	newWebhook := &dto.NewWebhook{}
	c.sev.Validate().Bind(gin, newWebhook)

	webhook, err := service.WebhookService().NewWebhook(newWebhook)
	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err))
		return
	}

	gin.JSON(200, webhook.ToDto())
}

func (c *WebhookController) GetName() string {
	return "webhook"
}

func (c *WebhookController) getEndpoint() string {
	return "/v1/webhooks"
}
