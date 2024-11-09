package main

import (
	"fmt"
	"log"
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
	INCREMENT_INTERVAL = time.Second * 10
	MOUSE_FILE         = "/dev/input/event19"
)

var (
	LeftClicks  uint = 0
	RightCLicks uint = 0
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
			fmt.Println("left click")
		}

		if ev.IsRightClick() {
			RightCLicks++
			fmt.Println("right click")
		}
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	sbClient := sb.CreateClient(db.LoadSecret())

	go TrackMouseInput()

	for {
		curSession, err := db.GetCurrentSession(sbClient)
		if err != nil {
			AddNewSession(sbClient)
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

		err = db.IncrementClickCount(sbClient, curSession, LeftClicks, RightCLicks)
		if err != nil {
			log.Fatal(err)
		}

		LeftClicks = 0
		RightCLicks = 0
	}
}
