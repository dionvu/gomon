package session

import (
	"time"

	"github.com/dionvu/gomon/hypr"
)

type Session struct {
	Id             string     `json:"id"`
	Start          time.Time  `json:"start"`
	End            time.Time  `json:"end"`
	Activity       []Activity `json:"activity"`
	LeftClicks     uint       `json:"left_clicks"`
	RightClicks    uint       `json:"right_clicks"`
	MiddleClicks   uint       `json:"middle_clicks"`
	XMovementMeter float64    `json:"mouse_movement_meter_x"`
	YMovementMeter float64    `json:"mouse_movement_meter_y"`
	KeyPresses     uint       `json:"key_presses"`
}

type Activity struct {
	Window       hypr.Window `json:"window"`
	TimeSpentMin float64     `json:"time_spent_min"`
}

// Returns Activity structs of all current windows,
// setting the time spent on them to 0 minutes.
func NewActivity(windows []hypr.Window) []Activity {
	activity := []Activity{}

	for _, win := range windows {
		activity = append(activity, Activity{
			Window: win,
		})
	}

	return activity
}

// Give the current windows, returns new Activity structs
// that are not already in the passed activity arr.
func FilterNewActivity(activity []Activity, windows []hypr.Window) []Activity {
	exists := make(map[hypr.Window]bool, len(activity))
	newActivity := []Activity{}

	for _, activ := range activity {
		exists[activ.Window] = true
	}

	for _, win := range windows {
		if !exists[win] {
			newActivity = append(newActivity, Activity{
				Window: win,
			})
		}

		exists[win] = true
	}

	return newActivity
}
