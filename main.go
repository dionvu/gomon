package main

import (
	"log"
	"math"
	"os"
	"runtime/debug"
	"time"

	"github.com/dionvu/temp/db"
	"github.com/dionvu/temp/event"
	"github.com/dionvu/temp/hypr"
	"github.com/dionvu/temp/session"
	sb "github.com/nedpals/supabase-go"
)

const (
	INCREMENT_INTERVAL = time.Minute
	MOUSE_FILE         = "/dev/input/event19"
	KEYBOARD_FILE      = "/dev/input/event16"
)

const HYPRCTL_UNIT_TO_METER = 0.0000244

var (
	LeftClicks      uint = 0
	RightCLicks     uint = 0
	MiddleClicks    uint = 0
	KeyboardPresses uint = 0
	MovementX       uint = 0
	MovementY       uint = 0
)

func AddNewSession(client *sb.Client) {
	windows, err := hypr.CurrentWindows()
	if err != nil {
		debug.PrintStack()
		log.Fatal(err)
	}

	activity := session.NewActivity(windows)

	_, err = db.AddHourSession(client, activity)
	if err != nil {
		log.Fatal(err)
	}
}

func HandleNewActivity(client *sb.Client, curSession session.Session, windows []hypr.Window) {
	newActivity := session.FilterNewActivity(curSession.Activity, windows)
	if len(newActivity) > 0 {
		err := db.UpdateNewActivity(client, curSession, newActivity)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TrackMouseInput() {
	f, err := os.Open(MOUSE_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b := make([]byte, event.SIZE)

	for {
		var ev event.InputEvent

		f.Read(b)

		ev, err = event.From(b)
		if err != nil {
			log.Fatal(err)
		}

		if ev.IsLeftClick() {
			LeftClicks++
		}

		if ev.IsRightClick() {
			RightCLicks++
		}

		if ev.IsMiddleClick() {
			MiddleClicks++
		}

		if ev.IsMouseMove() {
			switch ev.Code {
			case event.REL_X:
				MovementX += uint(math.Abs(float64(ev.Value)))
			case event.REL_Y:
				MovementY += uint(math.Abs(float64(ev.Value)))
			}
		}
	}
}

func TrackKeyboardInput() {
	f, err := os.Open(KEYBOARD_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b := make([]byte, event.SIZE)

	for {
		var ev event.InputEvent

		f.Read(b)

		ev, err = event.From(b)
		if err != nil {
			log.Fatal(err)
		}

		if ev.IsKeyboardPress() {
			KeyboardPresses++
		}
	}
}

func ResetCounters() {
	LeftClicks = 0
	RightCLicks = 0
	MiddleClicks = 0
	KeyboardPresses = 0
	MovementX = 0
	MovementY = 0
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	sbClient := sb.CreateClient(db.LoadSecret())

	go TrackMouseInput()
	go TrackKeyboardInput()

	for {
		var curSession session.Session

		curSession, err := db.GetCurrentSession(sbClient)
		if err != nil {
			AddNewSession(sbClient)
			curSession, err = db.GetCurrentSession(sbClient)
			if err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(INCREMENT_INTERVAL)

		windows, err := hypr.CurrentWindows()
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.IncrementActivityTime(sbClient, curSession, curSession.Activity, windows, INCREMENT_INTERVAL)
		if err != nil {
			log.Fatal(err)
		}

		HandleNewActivity(sbClient, curSession, windows)

		err = db.IncrementClickCount(sbClient, curSession, LeftClicks, RightCLicks, MiddleClicks)
		if err != nil {
			log.Fatal(err)
		}

		err = db.IncrementKeyboardPressCount(sbClient, curSession, KeyboardPresses)
		if err != nil {
			log.Fatal(err)
		}

		err = db.IncrementMouseMovement(sbClient, curSession, MovementX, MovementY)
		if err != nil {
			log.Fatal(err)
		}

		ResetCounters()
	}
}
