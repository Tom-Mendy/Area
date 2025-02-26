package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"area/repository"
	"area/schemas"
)

// Constructor

type TimerService interface {
	// Service interface functions
	GetServiceActionInfo() []schemas.Action
	GetServiceReactionInfo() []schemas.Reaction
	FindActionByName(name string) func(c chan string, option json.RawMessage, area schemas.Area)
	FindReactionByName(name string) func(option json.RawMessage, area schemas.Area) string
	// Service specific functions
	// Actions functions
	TimerActionSpecificHour(c chan string, option json.RawMessage, area schemas.Area)
	// Reactions functions
	TimerReactionGiveTime(option json.RawMessage, area schemas.Area) string
}

// timerService is a struct that provides services related to timers.
// It contains repositories for accessing timer, service, and area data,
// as well as information about the service itself.
//
// Fields:
// - repository: Interface for accessing timer data.
// - serviceRepository: Interface for accessing service data.
// - areaRepository: Interface for accessing area data.
// - serviceInfo: Information about the service.
type timerService struct {
	repository        repository.TimerRepository
	serviceRepository repository.ServiceRepository
	areaRepository    repository.AreaRepository
	serviceInfo       schemas.Service
}

// NewTimerService creates a new instance of TimerService with the provided repositories.
// It initializes the serviceInfo field with predefined values for the Timer service.
//
// Parameters:
//   - repository: an instance of TimerRepository for accessing timer data.
//   - serviceRepository: an instance of ServiceRepository for accessing service data.
//   - areaRepository: an instance of AreaRepository for accessing area data.
//
// Returns:
//   - TimerService: a new instance of TimerService.
func NewTimerService(
	repository repository.TimerRepository,
	serviceRepository repository.ServiceRepository,
	areaRepository repository.AreaRepository,
) TimerService {
	return &timerService{
		repository:        repository,
		serviceRepository: serviceRepository,
		areaRepository:    areaRepository,
		serviceInfo: schemas.Service{
			Name:        schemas.Timer,
			Description: "This service is a time service",
			Oauth:       false,
			Color:       "#BB00FF",
			Icon:        "https://api.iconify.design/mdi:clock.svg?color=%23FFFFFF",
		},
	}
}

// Service interface functions

// GetServiceInfo returns the service information.
// It retrieves the service information stored in the timerService instance.
//
// Returns:
//
//	schemas.Service: The service information.
func (service *timerService) GetServiceInfo() schemas.Service {
	return service.serviceInfo
}

// FindActionByName returns a function that matches the given action name.
// The returned function takes a channel, a JSON raw message, and an area schema as parameters.
// If the action name matches a specific time, it returns the TimerActionSpecificHour function.
// If no match is found, it returns nil.
//
// Parameters:
//   - name: The name of the action to find.
//
// Returns:
//   - A function that matches the given action name, or nil if no match is found.
func (service *timerService) FindActionByName(
	name string,
) func(c chan string, option json.RawMessage, area schemas.Area) {
	switch name {
	case string(schemas.SpecificTime):
		return service.TimerActionSpecificHour
	default:
		return nil
	}
}

// FindReactionByName returns a function that matches the given reaction name.
// The returned function takes a json.RawMessage option and a schemas.Area as parameters,
// and returns a string. If the reaction name does not match any known reactions,
// it returns nil.
//
// Parameters:
//   - name: The name of the reaction to find.
//
// Returns:
//   - A function that takes a json.RawMessage option and a schemas.Area, and returns a string.
//     If the reaction name does not match any known reactions, it returns nil.
func (service *timerService) FindReactionByName(
	name string,
) func(option json.RawMessage, area schemas.Area) string {
	switch name {
	case string(schemas.GiveTime):
		return service.TimerReactionGiveTime
	default:
		return nil
	}
}

// GetServiceActionInfo retrieves the service action information for the timer service.
// It initializes a default TimerActionSpecificHour value, marshals it into JSON, and
// updates the service information by finding the service by name. If any errors occur
// during marshaling or finding the service, they are printed to the console.
// The function returns a slice of Action containing the specific time action with the
// updated service information and default option.
//
// Returns:
//
//	[]schemas.Action: A slice containing the specific time action with the updated
//	service information and default option.
func (service *timerService) GetServiceActionInfo() []schemas.Action {
	defaultValue := schemas.TimerActionSpecificHour{
		Hour:   13,
		Minute: 7,
	}
	option, err := json.Marshal(defaultValue)
	if err != nil {
		println("error marshal timer option: " + err.Error())
	}
	service.serviceInfo, err = service.serviceRepository.FindByName(
		schemas.Timer,
	) // must update the serviceInfo
	if err != nil {
		println("error find service by name: " + err.Error())
	}
	return []schemas.Action{
		{
			Name:               string(schemas.SpecificTime),
			Description:        "This action is a specific time action",
			Service:            service.serviceInfo,
			Option:             option,
			MinimumRefreshRate: 10,
		},
	}
}

// GetServiceReactionInfo retrieves the reaction information for the timer service.
// It marshals a default value to JSON and updates the service information by finding
// the service by name. If any errors occur during these operations, they are printed
// to the console. The function returns a slice of Reaction structs containing the
// reaction details.
//
// Returns:
//
//	[]schemas.Reaction: A slice of Reaction structs with the reaction details.
func (service *timerService) GetServiceReactionInfo() []schemas.Reaction {
	defaultValue := struct{}{}
	option, err := json.Marshal(defaultValue)
	if err != nil {
		println("error marshal timer option: " + err.Error())
	}
	service.serviceInfo, err = service.serviceRepository.FindByName(
		schemas.Timer,
	) // must update the serviceInfo
	if err != nil {
		println("error find service by name: " + err.Error())
	}
	return []schemas.Reaction{
		{
			Name:        string(schemas.GiveTime),
			Description: "This reaction is a give time reaction",
			Service:     service.serviceInfo,
			Option:      option,
		},
	}
}

// Service specific functions

// getActualTime fetches the current time for the Europe/Paris timezone from the timeapi.io API.
// It returns a schemas.TimeApiResponse containing the time data or an error if the request fails.
//
// Returns:
//   - schemas.TimeApiResponse: The response containing the current time data.
//   - error: An error if the request creation, execution, or response decoding fails.
//
// Possible errors:
//   - schemas.ErrCreateRequest: If there is an error creating the HTTP request.
//   - schemas.ErrDoRequest: If there is an error executing the HTTP request.
//   - schemas.ErrDecode: If there is an error decoding the response body.
//   - fmt.Errorf: If the response status code is not 200 OK.
func getActualTime() (schemas.TimeApiResponse, error) {
	apiURL := "https://www.timeapi.io/api/time/current/zone?timeZone=Europe/Paris"

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return schemas.TimeApiResponse{}, schemas.ErrCreateRequest
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return schemas.TimeApiResponse{}, schemas.ErrDoRequest
	}

	if resp.StatusCode != http.StatusOK {
		return schemas.TimeApiResponse{}, fmt.Errorf("error status code %d", resp.StatusCode)
	}

	var result schemas.TimeApiResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return schemas.TimeApiResponse{}, schemas.ErrDecode
	}

	resp.Body.Close()
	return result, nil
}

// Actions functions

// TimerActionSpecificHour executes a timer action at a specific hour.
// It unmarshals the provided JSON option into a TimerActionSpecificHour struct,
// retrieves the current time from an external API, and updates the storage variable
// in the area repository if necessary. If the current time matches the specified hour
// and minute in the option, it sends a response message to the provided channel.
//
// Parameters:
//   - c: A channel to send the response message.
//   - option: A JSON raw message containing the timer action options.
//   - area: The area schema containing the storage variable.
//
// The function handles errors by printing error messages and sleeping for a second
// before returning. It also ensures that the storage variable is initialized and updated
// in the area repository if it is not already set.
func (service *timerService) TimerActionSpecificHour(
	c chan string,
	option json.RawMessage,
	area schemas.Area,
) {
	optionJSON := schemas.TimerActionSpecificHour{}

	err := json.Unmarshal(option, &optionJSON)
	if err != nil {
		println("error unmarshal timer option: " + err.Error())
		time.Sleep(time.Second)
		return
	}

	actualTimeApi, err := getActualTime()
	if err != nil {
		println("error get actual time" + err.Error())
		time.Sleep(time.Second)
		return
	}

	databaseStored := schemas.TimerActionSpecificHourStorage{}
	err = json.Unmarshal(area.StorageVariable, &databaseStored)
	if err != nil {
		toto := struct{}{}
		err = json.Unmarshal(area.StorageVariable, &toto)
		if err != nil {
			println("error unmarshalling storage variable: " + err.Error())
			return
		} else {
			println("initializing storage variable")
			databaseStored = schemas.TimerActionSpecificHourStorage{
				Time: time.Now(),
			}
			area.StorageVariable, err = json.Marshal(databaseStored)
			if err != nil {
				println("error marshalling storage variable: " + err.Error())
				return
			}
			err = service.areaRepository.Update(area)
			if err != nil {
				println("error updating area: " + err.Error())
				return
			}
		}
	}

	if databaseStored.Time.IsZero() {
		println("initializing storage variable")
		databaseStored = schemas.TimerActionSpecificHourStorage{
			Time: time.Now(),
		}
		area.StorageVariable, err = json.Marshal(databaseStored)
		if err != nil {
			println("error marshalling storage variable: " + err.Error())
			return
		}
		err = service.areaRepository.Update(area)
		if err != nil {
			println("error updating area: " + err.Error())
			return
		}
	}

	// generate time.Time from actualTimeApi
	actualTime := time.Date(
		actualTimeApi.Year,
		time.Month(actualTimeApi.Month),
		actualTimeApi.Day,
		actualTimeApi.Hour,
		actualTimeApi.Minute,
		actualTimeApi.Seconds,
		actualTimeApi.MilliSeconds,
		time.Local,
	)

	if databaseStored.Time.Before(actualTime) {
		if actualTime.Hour() == optionJSON.Hour && actualTimeApi.Minute == optionJSON.Minute {
			response := "current time is " + actualTimeApi.Time
			databaseStored.Time = time.Now().Add(time.Minute)
			area.StorageVariable, err = json.Marshal(databaseStored)
			if err != nil {
				println("error marshalling storage variable: " + err.Error())
				return
			}
			err = service.areaRepository.Update(area)
			if err != nil {
				println("error updating area: " + err.Error())
				return
			}
			println(response)
			c <- response
		}
	}

	if (area.Action.MinimumRefreshRate) > area.ActionRefreshRate {
		time.Sleep(time.Second * time.Duration(area.Action.MinimumRefreshRate))
	} else {
		time.Sleep(time.Second * time.Duration(area.ActionRefreshRate))
	}
}

// Reactions functions

// TimerReactionGiveTime retrieves the current time from an external API and returns it as a string.
// If there is an error while fetching the time, it logs the error and returns an error message.
//
// Parameters:
//
//	option - a JSON raw message containing additional options (currently unused).
//	area - a schemas.Area object representing the area (currently unused).
//
// Returns:
//
//	A string containing the current time or an error message if the time could not be retrieved.
func (service *timerService) TimerReactionGiveTime(
	option json.RawMessage,
	area schemas.Area,
) string {
	actualTimeApi, err := getActualTime()
	if err != nil {
		println("error get actual time" + err.Error())
		return "error get actual time"
	} else {
		response := "current time is " + actualTimeApi.Time
		println(response)
		return response
	}
}
