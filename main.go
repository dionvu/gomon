package main

import (
	"log"
	"os"
	"runtime/debug"
	"time"

	"github.com/dionvu/gomon/db"
	"github.com/dionvu/gomon/hypr"
	"github.com/dionvu/gomon/input"
	"github.com/dionvu/gomon/session"
	sb "github.com/nedpals/supabase-go"
)

const (
	INCREMENT_INTERVAL = time.Minute
	MOUSE_FILE         = "/dev/input/event19"
	KEYBOARD_FILE      = "/dev/input/event16"
)

var (
	LeftClicks      uint    = 0
	RightCLicks     uint    = 0
	MiddleClicks    uint    = 0
	KeyboardPresses uint    = 0
	XMovementMeter  float64 = 0
	YMovementMeter  float64 = 0
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

	b := make([]byte, input.SIZE)

	for {
		var ev input.Event

		f.Read(b)

		ev, err = input.From(b)
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
	}
}

func TrackMouseMovement() {
	f, err := os.Open(MOUSE_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b := make([]byte, input.SIZE)

	for {
		var ev input.Event

		f.Read(b)

		ev, err = input.From(b)
		if err != nil {
			log.Fatal(err)
		}

		if ev.IsMouseMove() {
			switch ev.Code {
			case input.REL_X:
				XMovementMeter += ev.Value.Abs().Meter()
			case input.REL_Y:
				YMovementMeter += ev.Value.Abs().Meter()
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

	b := make([]byte, input.SIZE)

	for {
		var ev input.Event

		f.Read(b)

		ev, err = input.From(b)
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
	XMovementMeter = 0
	YMovementMeter = 0
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	sbClient := sb.CreateClient(db.LoadSecret())

	ses, err := db.GetCurrentSession(sbClient)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.ArchieveSession(sbClient, ses)
	if err != nil {
		log.Fatal(err)
	}

	// go TrackMouseInput()
	// go TrackKeyboardInput()
	// go TrackMouseMovement()
	//
	// for {
	// 	var curSession session.Session
	//
	// 	curSession, err := db.GetCurrentSession(sbClient)
	// 	if err != nil {
	// 		AddNewSession(sbClient)
	// 		curSession, err = db.GetCurrentSession(sbClient)
	// 		if err != nil {
	// 			log.Fatal(err)
	// 		}
	// 	}
	//
	// 	time.Sleep(INCREMENT_INTERVAL)
	//
	// 	windows, err := hypr.CurrentWindows()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	_, err = db.IncrementActivityTime(sbClient, curSession, curSession.Activity, windows, INCREMENT_INTERVAL)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	HandleNewActivity(sbClient, curSession, windows)
	//
	// 	err = db.IncrementClickCount(sbClient, curSession, LeftClicks, RightCLicks, MiddleClicks)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	err = db.IncrementKeyboardPressCount(sbClient, curSession, KeyboardPresses)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	err = db.IncrementMouseMovement(sbClient, curSession, XMovementMeter, YMovementMeter)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	ResetCounters()
	// }
}
