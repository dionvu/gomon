package session

import (
	"time"

	"github.com/dionvu/temp/hypr"
)

type Session struct {
	Id       string     `json:"id"`
	Start    time.Time  `json:"start"`
	End      time.Time  `json:"end"`
	Activity []Activity `json:"activity"`
}

type Activity struct {
	Window       hypr.Window `json:"window"`
	TimeSpentMin int         `json:"time_spent_min"`
}

func NewActivity() ([]Activity, error) {
	activity := []Activity{}

	windows, err := hypr.Windows()
	if err != nil {
		return activity, err
	}

	for _, window := range windows {
		activity = append(activity, Activity{
			Window:       window,
			TimeSpentMin: 0,
		})
	}

	return activity, nil
}
