package schemas

import "errors"

type OpenweathermapAction string

const (
	SpecificWeather     OpenweathermapAction = "SpecificWeather"
	SpecificTemperature OpenweathermapAction = "SpecificTemperature"
)

type OpenweathermapReaction string

const (
	CurrentWeather     OpenweathermapReaction = "CurrentWeather"
	CurrentTemperature OpenweathermapReaction = "CurrentTemperature"
)

// https://openweathermap.org/weather-conditions

type WeatherCondition string

const (
	// Thunderstorm
	Thunderstorm WeatherCondition = "Thunderstorm"
	// Drizzle
	Drizzle WeatherCondition = "Drizzle"
	// Rain
	Rain WeatherCondition = "Rain"
	// Snow
	Snow WeatherCondition = "Snow"
	// Atmosphere
	Mist    WeatherCondition = "Mist"
	Smoke   WeatherCondition = "Smoke"
	Haze    WeatherCondition = "Haze"
	Dust    WeatherCondition = "Dust"
	Fog     WeatherCondition = "Fog"
	Sand    WeatherCondition = "Sand"
	Ash     WeatherCondition = "Ash"
	Squall  WeatherCondition = "Squall"
	Tornado WeatherCondition = "Tornado"
	// Clear
	Clear WeatherCondition = "Clear"
	// Clouds
	Clouds WeatherCondition = "Clouds"
)

type OpenweathermapActionSpecificWeather struct {
	City    string           `json:"city"`
	Weather WeatherCondition `json:"weather"`
}

type OpenweathermapActionSpecificTemperature struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature"`
}

// all reaction options schema
type OpenweathermapReactionOption struct {
	City string `json:"city"`
}

type OpenweathermapCityCoordinatesResponse struct {
	Name    string  `json:"name"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	State   string  `json:"state"`
}

type OpenweathermapCoordinatesWeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int              `json:"id"`
		Main        WeatherCondition `json:"main"`
		Description string           `json:"description"`
		Icon        string           `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Rain struct {
		OneH float64 `json:"1h"`
	} `json:"rain"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int    `json:"type"`
		ID      int    `json:"id"`
		Country string `json:"country"`
		Sunrise int    `json:"sunrise"`
		Sunset  int    `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

type OpenweathermapReactionGiveTime struct{}

type OpenweathermapReactionApiResponse struct{}

var ErrOpenWeatherMapApiKeyNotSet = errors.New("OPENWEATHERMAP_API_KEY is not set")
