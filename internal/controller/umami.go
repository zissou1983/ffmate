package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/internal/dto"
	"github.com/welovemedia/ffmate/sev"
)

type UmamiController struct {
	sev.Controller
	sev *sev.Sev

	Prefix string
}

func (c *UmamiController) Setup(s *sev.Sev) {
	c.sev = s
	s.Gin().POST(c.getEndpoint(), c.analytics)
}

func (c *UmamiController) analytics(gin *gin.Context) {
	a := &dto.Umami{}
	c.sev.Validate().Bind(gin, a)
	c.sev.Metrics().GaugeVec("umami").WithLabelValues(a.Payload.Url, a.Payload.Screen, a.Payload.Langugage).Inc()
	gin.Status(204)
}

func (c *UmamiController) GetName() string {
	return "umami"
}

func (c *UmamiController) getEndpoint() string {
	return "/metrics/umami"
}
