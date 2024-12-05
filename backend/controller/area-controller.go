package controller

import (
	"github.com/gin-gonic/gin"

	"area/service"
)

type AreaController interface {
	CreateArea(ctx *gin.Context) (string, error)
}

type areaController struct {
	service service.AreaService
}

func NewAreaController(service service.AreaService) AreaController {
	return &areaController{
		service: service,
	}
}

func (controller *areaController) CreateArea(ctx *gin.Context) (string, error) {
	if ctx == nil {
		return "", nil
	}
	return "status: success", nil
	// return controller.service.CreateArea()
}
