package session

import (
	"log"
	"os"
	"time"

	"github.com/dionvu/gomon/input"
)

type Session struct {
	Id           string    `json:"id"`
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	LeftClicks   uint      `json:"left_clicks"`
	RightClicks  uint      `json:"right_clicks"`
	MiddleClicks uint      `json:"middle_clicks"`
	XMovement    uint      `json:"x_mouse_movement"`
	YMovement    uint      `json:"y_mouse_movement"`
	KeyPresses   uint      `json:"key_presses"`
}

type Tracker struct {
	mouseFile       *os.File
	keyboardFile    *os.File
	trackpadFile    *os.File
	LeftClicks      uint
	RightClicks     uint
	MiddleClicks    uint
	KeyboardPresses uint
	XMovement       uint
	YMovement       uint
}

func NewTracker(keyboardEvent string, mouseEvent string, trackpadEvent string) Tracker {
	tracker := Tracker{
		LeftClicks:      0,
		RightClicks:     0,
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

	f, err = os.Open(trackpadEvent)
	if err != nil {
		log.Fatal(err)
	}

	tracker.trackpadFile = f

	return tracker
}

func (t *Tracker) Reset() {
	t.LeftClicks = 0
	t.RightClicks = 0
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
			tracker.RightClicks++
		}

		if ev.IsMiddleClick() {
			tracker.MiddleClicks++
		}
	}
}

func (tracker *Tracker) TrackMouseMovement() {
	b := make([]byte, input.SIZE)

	for {
		tracker.mouseFile.Read(b)

		event, err := input.From(b)
		if err != nil {
			log.Fatal(err)
		}

		if event.IsMouseMove() {
			if event.Code.IsRelX() {
				tracker.IncrementX(event.Value)
			}

			if event.Code.IsRelY() {
				tracker.IncrementY(event.Value)
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

func (tracker *Tracker) TrackTrackpadInput() {
	b := make([]byte, input.SIZE)
	for {
		tracker.trackpadFile.Read(b)
		event, err := input.From(b)
		if err != nil {
			log.Fatal(err)
		}
		if event.IsTrackpadRightClick() {
			tracker.RightClicks++
			tracker.LeftClicks-- // For some reason also adds a left
		}
		if event.IsTrackpadLeftClick() {
			tracker.LeftClicks++
		}
	}
}

func (tracker *Tracker) IncrementX(x input.EventValue) {
	tracker.XMovement += uint(x.Abs())
}

func (tracker *Tracker) IncrementY(y input.EventValue) {
	tracker.YMovement += uint(y.Abs())
}

func (tracker *Tracker) ListenAll() {
	go tracker.TrackKeyboardInput()
	go tracker.TrackMouseInput()
	go tracker.TrackMouseMovement()
	go tracker.TrackTrackpadInput()
}

func (tracker *Tracker) Close() {
	tracker.mouseFile.Close()
	tracker.keyboardFile.Close()
}
