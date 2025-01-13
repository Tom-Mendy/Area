package controller

import (
	"fmt"

	"area/schemas"
	"area/service"
)

type ActionController interface {
	GetActionsInfo(id uint64) (response []schemas.Action, err error)
	GetServiceInfo(id uint64) (response schemas.Service, err error)
}

type actionController struct {
	service service.ActionService
}

func NewActionController(service service.ActionService) ActionController {
	return &actionController{
		service: service,
	}
}

func (controller *actionController) GetActionsInfo(
	id uint64,
) (response []schemas.Action, err error) {
	response, err = controller.service.GetActionsInfo(id)
	if err != nil {
		return nil, fmt.Errorf("unable to get actions info because %w", err)
	}
	return response, nil
}

func (controller *actionController) GetServiceInfo(
	id uint64,
) (response schemas.Service, err error) {
	action, err := controller.service.FindById(id)
	if err != nil {
		return response, fmt.Errorf("unable to get actions info because %w", err)
	}
	response = action.Service
	return response, nil
}
