package session

import (
	"fmt"
	"time"

	"github.com/dionvu/temp/hypr"
	"github.com/google/uuid"
)

type Session struct {
	Id       string     `json:"id"`
	Start    time.Time  `json:"start"`
	End      time.Time  `json:"end"`
	Activity []Activity `json:"activity"`
}

type Activity struct {
	Window       hypr.Window `json:"window"`
	TimeSpentMin float64     `json:"time_spent_min"`
}

// Returns activity structs of all current windows,
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

// Checks if session has a valid uuid.
func (s Session) IsValid() bool {
	fmt.Print(s.Id)
	_, err := uuid.Parse(s.Id)
	return err == nil
}
