package controller

import (
	"area/service"
)

type TimerController interface{}

type timerController struct {
	service service.TimerService
}

func NewTimerController(service service.TimerService) TimerController {
	return &timerController{
		service: service,
	}
}
