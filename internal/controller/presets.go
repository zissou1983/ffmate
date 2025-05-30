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

type PresetController struct {
	sev.Controller
	sev *sev.Sev

	Prefix string
}

func (c *PresetController) Setup(s *sev.Sev) {
	c.sev = s
	s.Gin().DELETE(c.Prefix+c.getEndpoint()+"/:uuid", c.deletePreset)
	s.Gin().POST(c.Prefix+c.getEndpoint(), c.addPreset)
	s.Gin().PUT(c.Prefix+c.getEndpoint()+"/:uuid", c.updatePreset)
	s.Gin().GET(c.Prefix+c.getEndpoint(), interceptor.PageLimit, c.listPresets)
	s.Gin().GET(c.Prefix+c.getEndpoint()+"/:uuid", c.getPreset)
}

// @Summary Delete a preset
// @Description Delete a preset by its uuid
// @Tags presets
// @Param uuid path string true "the presets uuid"
// @Produce json
// @Success 204
// @Router /presets/{uuid} [delete]
func (c *PresetController) deletePreset(gin *gin.Context) {
	uuid := gin.Param("uuid")
	err := service.PresetService().DeletePreset(uuid)

	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err, "https://docs.ffmate.io/docs/presets#deleting-a-preset"))
		return
	}

	gin.AbortWithStatus(204)
}

// @Summary List all presets
// @Description	List all existing presets
// @Tags presets
// @Produce json
// @Success 200 {object} []dto.Preset
// @Router /presets [get]
func (c *PresetController) listPresets(gin *gin.Context) {
	presets, total, err := service.PresetService().ListPresets(gin.GetInt("page"), gin.GetInt("perPage"))
	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err, "https://docs.ffmate.io/docs/presets#listing-presets"))
		return
	}

	gin.Header("X-Total", fmt.Sprintf("%d", total))

	// Transform each preset to its DTO
	var presetDTOs = []dto.Preset{}
	for _, preset := range *presets {
		presetDTOs = append(presetDTOs, *preset.ToDto())
	}

	gin.JSON(200, presetDTOs)
}

// @Summary Add a new preset
// @Description	Add a new preset
// @Tags presets
// @Accept json
// @Param request body dto.NewPreset true "new preset"
// @Produce json
// @Success 200 {object} dto.Preset
// @Router /presets [post]
func (c *PresetController) addPreset(gin *gin.Context) {
	newPreset := &dto.NewPreset{}
	if !c.sev.Validate().Bind(gin, newPreset) {
		return
	}

	preset, err := service.PresetService().NewPreset(newPreset)
	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err, "https://docs.ffmate.io/docs/presets#creating-a-preset"))
		return
	}

	gin.JSON(200, preset.ToDto())
}

// @Summary Get a preset
// @Description	Get a preset
// @Tags presets
// @Produce json
// @Success 200 {object} dto.Preset
// @Router /presets/{uuid} [get]
func (c *PresetController) getPreset(gin *gin.Context) {
	presetUuid := gin.Param("uuid")

	preset, err := service.PresetService().FindByUuid(presetUuid)
	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err, "https://docs.ffmate.io/docs/presets#getting-a-single-preset"))
		return
	}

	gin.JSON(200, preset.ToDto())
}

// @Summary Update a preset
// @Description	Update a preset
// @Tags presets
// @Accept json
// @Param request body dto.NewPreset true "new preset"
// @Produce json
// @Success 200 {object} dto.Preset
// @Router /presets [put]
func (c *PresetController) updatePreset(gin *gin.Context) {
	uuid := gin.Param("uuid")
	newPreset := &dto.NewPreset{}
	if !c.sev.Validate().Bind(gin, newPreset) {
		return
	}

	preset, err := service.PresetService().UpdatePreset(uuid, newPreset)
	if err != nil {
		gin.JSON(400, exceptions.HttpBadRequest(err, "https://docs.ffmate.io/docs/presets#updating-a-preset"))
		return
	}

	gin.JSON(200, preset.ToDto())
}

func (c *PresetController) GetName() string {
	return "preset"
}

func (c *PresetController) getEndpoint() string {
	return "/v1/presets"
}
