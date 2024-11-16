package session

import (
	"log"
	"os"
	"time"

	"github.com/dionvu/gomon/hypr"
	"github.com/dionvu/gomon/input"
)

type Session struct {
	Id           string     `json:"id"`
	Start        time.Time  `json:"start"`
	End          time.Time  `json:"end"`
	Activity     []Activity `json:"activity"`
	LeftClicks   uint       `json:"left_clicks"`
	RightClicks  uint       `json:"right_clicks"`
	MiddleClicks uint       `json:"middle_clicks"`
	XMovement    uint       `json:"x_mouse_movement"`
	YMovement    uint       `json:"y_mouse_movement"`
	KeyPresses   uint       `json:"key_presses"`
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

type Tracker struct {
	mouseFile       *os.File
	keyboardFile    *os.File
	LeftClicks      uint
	RightCLicks     uint
	MiddleClicks    uint
	KeyboardPresses uint
	XMovement       uint
	YMovement       uint
}

func NewTracker(keyboardEvent string, mouseEvent string) Tracker {
	tracker := Tracker{
		LeftClicks:      0,
		RightCLicks:     0,
		MiddleClicks:    0,
		KeyboardPresses: 0,
		XMovement:       0,
		YMovement:       0,
	}

	f, err := os.Open(keyboardEvent)
	if err != nil {
		log.Fatal(err)
	}

	tracker.keyboardFile = f

	f, err = os.Open(mouseEvent)
	if err != nil {
		log.Fatal(err)
	}

	tracker.mouseFile = f

	return tracker
}

func (t *Tracker) Reset() {
	t.LeftClicks = 0
	t.RightCLicks = 0
	t.MiddleClicks = 0
	t.KeyboardPresses = 0
	t.XMovement = 0
	t.YMovement = 0
}

func (tracker *Tracker) TrackMouseInput() {
	b := make([]byte, input.SIZE)

	for {
		var ev input.Event

		tracker.mouseFile.Read(b)

		ev, err := input.From(b)
		if err != nil {
			log.Fatal(err)
		}

		if ev.IsLeftClick() {
			tracker.LeftClicks++
		}

		if ev.IsRightClick() {
			tracker.RightCLicks++
		}

		if ev.IsMiddleClick() {
			tracker.MiddleClicks++
		}
	}
}

func (tracker *Tracker) TrackMouseMovement() {
	b := make([]byte, input.SIZE)

	for {
		var ev input.Event

		tracker.mouseFile.Read(b)

		event, err := input.From(b)
		if err != nil {
			log.Fatal(err)
		}

		if event.IsMouseMove() {
			if event.Code.IsRelX() {
				tracker.IncrementX(ev.Value)
			}

			if event.Code.IsRelY() {
				tracker.IncrementY(ev.Value)
			}
		}
	}
}

func (tracker *Tracker) TrackKeyboardInput() {
	b := make([]byte, input.SIZE)

	for {
		var event input.Event

		tracker.keyboardFile.Read(b)

		event, err := input.From(b)
		if err != nil {
			log.Fatal(err)
		}

		if event.IsKeyboardPress() {
			tracker.KeyboardPresses++
		}
	}
}

func (tracker *Tracker) IncrementX(x input.EventValue) {
	tracker.XMovement += uint(x.Abs())
}

func (tracker *Tracker) IncrementY(y input.EventValue) {
	tracker.YMovement += uint(y.Abs())
}

// Listens to all events async.
func (tracker *Tracker) ListenAll() {
	go tracker.TrackKeyboardInput()
	go tracker.TrackMouseInput()
	go tracker.TrackMouseMovement()
}

func (tracker *Tracker) Close() {
	tracker.mouseFile.Close()
	tracker.keyboardFile.Close()
}
