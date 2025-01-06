package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/welovemedia/ffmate/pkg/dto"
	"github.com/welovemedia/ffmate/sev"
)

type VersionController struct {
	sev.Controller
	sev    *sev.Sev
	Prefix string
}

func (v *VersionController) Setup(s *sev.Sev) {
	v.sev = s
	s.Gin().GET(v.Prefix+v.getEndpoint(), v.getVersion)
}

//	@Summary		Get ffmate version
//	@Description	Get ffmate version
//	@Tags			version
//	@Produce		json
//	@Success		200	{object}	dto.Version
//	@Router			/version [get]
func (v *VersionController) getVersion(gin *gin.Context) {
	gin.JSON(200, &dto.Version{Version: v.sev.AppVersion()})
}

func (v *VersionController) GetName() string {
	return "version"
}

func (v *VersionController) getEndpoint() string {
	return "/v1/version"
}
