package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"area/repository"
	"area/schemas"
)

// Constructor

type OpenweathermapService interface {
	// Service interface functions
	GetServiceActionInfo() []schemas.Action
	GetServiceReactionInfo() []schemas.Reaction
	FindActionbyName(name string) func(c chan string, option string, idArea uint64)
	FindReactionbyName(name string) func(option string, idArea uint64) string
	GetActionsName() []string
	GetReactionsName() []string
	// Service specific functions
	// Actions functions
	OpenweathermapActionSpecificWeather(
		c chan string,
		option string,
		idArea uint64,
	)
	// Reactions functions
	OpenweathermapReactionCurrentWeather(
		option string,
		idArea uint64,
	) string
	OpenweathermapActionSpecificTemperature(
		c chan string,
		option string,
		idArea uint64,
	)
	OpenweathermapReactionCurrentTemperature(
		option string,
		idArea uint64,
	) string
}

type openweathermapService struct {
	repository        repository.OpenweathermapRepository
	serviceRepository repository.ServiceRepository
	actionsName       []string
	reactionsName     []string
	serviceInfo       schemas.Service
}

func NewOpenweathermapService(
	repository repository.OpenweathermapRepository,
	serviceRepository repository.ServiceRepository,
) OpenweathermapService {
	return &openweathermapService{
		repository:        repository,
		serviceRepository: serviceRepository,
		serviceInfo: schemas.Service{
			Name:        schemas.Openweathermap,
			Description: "This service is a weather service",
		},
	}
}

// Service interface functions

func (service *openweathermapService) GetServiceInfo() schemas.Service {
	return service.serviceInfo
}

func (service *openweathermapService) FindActionbyName(
	name string,
) func(c chan string, option string, idArea uint64) {
	switch name {
	case string(schemas.SpecificWeather):
		return service.OpenweathermapActionSpecificWeather
	case string(schemas.SpecificTemperature):
		return service.OpenweathermapActionSpecificTemperature
	default:
		return nil
	}
}

func (service *openweathermapService) FindReactionbyName(
	name string,
) func(option string, idArea uint64) string {
	switch name {
	case string(schemas.CurrentWeather):
		return service.OpenweathermapReactionCurrentWeather
	default:
		return nil
	}
}

func (service *openweathermapService) GetServiceActionInfo() []schemas.Action {
	service.actionsName = append(
		service.actionsName,
		string(schemas.SpecificWeather),
		string(schemas.SpecificTemperature),
	)
	return []schemas.Action{
		{
			Name:        string(schemas.SpecificWeather),
			Description: "This action is a specific weather action",
			Service:     service.serviceRepository.FindByName(schemas.Openweathermap),
			Option:      "{\"city\": \"\", \"weather\": \"\"}",
		},
		{
			Name:        string(schemas.SpecificTemperature),
			Description: "This action is a specific temperature action",
			Service:     service.serviceRepository.FindByName(schemas.Openweathermap),
			Option:      "{\"city\": \"\", \"temperature\": 0}",
		},
	}
}

func (service *openweathermapService) GetServiceReactionInfo() []schemas.Reaction {
	service.reactionsName = append(
		service.reactionsName,
		string(schemas.CurrentWeather),
		string(schemas.CurrentTemperature),
	)
	return []schemas.Reaction{
		{
			Name:        string(schemas.CurrentWeather),
			Description: "This reaction is a current weather reaction",
			Service:     service.serviceRepository.FindByName(schemas.Openweathermap),
			Option:      "{\"city\": \"\"}",
		},
		{
			Name:        string(schemas.CurrentTemperature),
			Description: "This reaction is a current teamperature reaction",
			Service:     service.serviceRepository.FindByName(schemas.Openweathermap),
			Option:      "{\"city\": \"\"}",
		},
	}
}

func (service *openweathermapService) GetActionsName() []string {
	return service.actionsName
}

func (service *openweathermapService) GetReactionsName() []string {
	return service.reactionsName
}

// Service specific functions

func getCoordinatesOfCity(city string) (coordinates struct {
	Lat float64
	Lon float64
}, err error,
) {
	APIKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if APIKey == "" {
		return coordinates, schemas.ErrOpenWeatherMapApiKeyNotSet
	}
	apiURL := "http://api.openweathermap.org/geo/1.0/direct"
	data := url.Values{}
	data.Set("q", city)
	data.Set("limit", "1")
	data.Set("appid", APIKey)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return coordinates, fmt.Errorf("unable to create request because %w", err)
	}

	req.URL.RawQuery = data.Encode()
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return coordinates, fmt.Errorf("unable to make request because %w", err)
	}

	var result []schemas.OpenweathermapCityCoordinatesResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return coordinates, fmt.Errorf(
			"unable to decode response because %w",
			err,
		)
	}
	coordinates.Lat = result[0].Lat
	coordinates.Lon = result[0].Lon
	return coordinates, nil
}

func getWeatherOfCoodinate(coordinates struct {
	Lat float64
	Lon float64
},
) (weather schemas.OpenweathermapCoordinatesWeatherResponse, err error) {
	APIKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if APIKey == "" {
		return weather, schemas.ErrOpenWeatherMapApiKeyNotSet
	}
	apiURL := "https://api.openweathermap.org/data/2.5/weather"
	data := url.Values{}
	data.Set("lat", fmt.Sprintf("%f", coordinates.Lat))
	data.Set("lon", fmt.Sprintf("%f", coordinates.Lon))
	data.Set("appid", APIKey)
	data.Set("units", "metric") // to get temperature in celsius

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return weather, fmt.Errorf("unable to create request because %w", err)
	}

	req.URL.RawQuery = data.Encode()
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return weather, fmt.Errorf("unable to make request because %w", err)
	}

	var result schemas.OpenweathermapCoordinatesWeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return weather, fmt.Errorf(
			"unable to decode response because %w",
			err,
		)
	}
	weather = result
	return weather, nil
}

// Actions functions

func (service *openweathermapService) OpenweathermapActionSpecificWeather(
	c chan string,
	option string,
	idArea uint64,
) {
	optionJSON := schemas.OpenweathermapActionSpecificWeather{}

	err := json.Unmarshal([]byte(option), &optionJSON)
	if err != nil {
		println("error unmarshal weather option: " + err.Error())
		time.Sleep(time.Second)
		return
	}

	coordinates, err := getCoordinatesOfCity(optionJSON.City)
	if err != nil {
		fmt.Println(err)
	}
	weatherOfSpecifiedCity, err := getWeatherOfCoodinate(coordinates)
	if err != nil {
		println("error get actual weather info" + err.Error())
	} else {
		if weatherOfSpecifiedCity.Weather[0].Main == optionJSON.Weather {
			response := "current weather in " + optionJSON.City + " is " + string(weatherOfSpecifiedCity.Weather[0].Main)
			println(response)
			c <- response
		}
	}
	time.Sleep(time.Minute)
}

func (service *openweathermapService) OpenweathermapActionSpecificTemperature(
	c chan string,
	option string,
	idArea uint64,
) {
	optionJSON := schemas.OpenweathermapActionSpecificTemperature{}

	err := json.Unmarshal([]byte(option), &optionJSON)
	if err != nil {
		println("error unmarshal temperature option: " + err.Error())
		time.Sleep(time.Second)
		return
	}

	coordinates, err := getCoordinatesOfCity(optionJSON.City)
	if err != nil {
		fmt.Println(err)
	}
	weatherOfSpecifiedCity, err := getWeatherOfCoodinate(coordinates)

	if err != nil {
		println("error get actual temperature info" + err.Error())
	} else {
		if weatherOfSpecifiedCity.Main.Temp == optionJSON.Temperature {
			response := "current temperature in " + optionJSON.City + " is " + fmt.Sprintf("%f", weatherOfSpecifiedCity.Main.Temp) + "°C"
			println(response)
			c <- response
		}
	}
	time.Sleep(time.Minute)
}

// Reactions functions

func (service *openweathermapService) OpenweathermapReactionCurrentWeather(
	option string,
	idArea uint64,
) string {
	optionJSON := schemas.OpenweathermapReactionOption{}

	err := json.Unmarshal([]byte(option), &optionJSON)
	if err != nil {
		println("error unmarshal weather option: " + err.Error())
		time.Sleep(time.Second)
		return "error unmarshal weather option: " + err.Error()
	}
	coordinates, err := getCoordinatesOfCity(optionJSON.City)
	if err != nil {
		fmt.Println(err)
	}
	weatherOfSpecifiedCity, err := getWeatherOfCoodinate(coordinates)
	if err != nil {
		println("error get actual weather info" + err.Error())
		return "error get actual weather info" + err.Error()
	} else {
		response := "current weather in " + optionJSON.City + " is " + string(weatherOfSpecifiedCity.Weather[0].Main)
		println(response)
		return response
	}
}

func (service *openweathermapService) OpenweathermapReactionCurrentTemperature(
	option string,
	idArea uint64,
) string {
	optionJSON := schemas.OpenweathermapReactionOption{}

	err := json.Unmarshal([]byte(option), &optionJSON)
	if err != nil {
		println("error unmarshal temperature option: " + err.Error())
		time.Sleep(time.Second)
		return "error unmarshal temperature option: " + err.Error()
	}
	coordinates, err := getCoordinatesOfCity(optionJSON.City)
	if err != nil {
		fmt.Println(err)
	}
	weatherOfSpecifiedCity, err := getWeatherOfCoodinate(coordinates)
	if err != nil {
		println("error get actual temperature info" + err.Error())
		return "error get actual temperature info" + err.Error()
	} else {
		response := "current temperature in " + optionJSON.City + " is " + fmt.Sprintf("%f", weatherOfSpecifiedCity.Main.Temp) + "°C"
		println(response)
		return response
		// TODO: save to database
	}
}
