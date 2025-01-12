package schemas

import "time"

type TimerAction string

const (
	SpecificTime TimerAction = "SpecificTime"
)

type TimerReaction string

const (
	GiveTime TimerReaction = "GiveTime"
)

type TimerActionSpecificHour struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

type TimerActionSpecificHourStorage struct {
	Time time.Time `json:"time"`
}

type TimerReactionGiveTime struct{}

type TimeApiResponse struct {
	Year         int    `json:"year"`
	Month        int    `json:"month"`
	Day          int    `json:"day"`
	Hour         int    `json:"hour"`
	Minute       int    `json:"minute"`
	Seconds      int    `json:"seconds"`
	MilliSeconds int    `json:"milliSeconds"`
	DateTime     string `json:"dateTime"`
	Date         string `json:"date"`
	Time         string `json:"time"`
	TimeZone     string `json:"timeZone"`
	DayOfWeek    string `json:"dayOfWeek"`
	DstActive    bool   `json:"dstActive"`
}
