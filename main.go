package main

import (
	"log"
	"time"

	"github.com/dionvu/gomon/db"
	"github.com/dionvu/gomon/hypr"
	"github.com/dionvu/gomon/session"
	sb "github.com/nedpals/supabase-go"
)

const (
	INCREMENT_INTERVAL = time.Second * 2
	KEYBOARD_FILE      = "/dev/input/event16"
	MOUSE_FILE         = "/dev/input/event19"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	sbClient := sb.CreateClient(db.LoadSecret())

	tracker := session.NewTracker(KEYBOARD_FILE, MOUSE_FILE)
	defer tracker.Close()

	tracker.ListenAll()

	for {
		var curSession session.Session
		var curWindows []hypr.Window

		curSession, err := db.GetCurrentSession(sbClient)
		if err != nil {
			curWindows, err = hypr.CurrentWindows()
			if err != nil {
				log.Fatal(err)
			}

			activity := session.NewActivity(curWindows)

			_, err = db.AddSession(sbClient, activity)
			if err != nil {
				log.Fatal(err)
			}

			curSession, err = db.GetCurrentSession(sbClient)
			if err != nil {
				log.Fatal(err)
			}
		}

		time.Sleep(INCREMENT_INTERVAL)

		_, err = db.ArchivePastSessions(sbClient)
		if err != nil {
			log.Fatal(err)
		}

		db.DropSessions(sbClient, time.Now().Add(-24*time.Hour))

		curWindows, err = hypr.CurrentWindows()
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.IncrementActivityTime(sbClient, curSession, curSession.Activity, curWindows, INCREMENT_INTERVAL)
		if err != nil {
			log.Fatal(err)
		}

		newActivity := session.FilterNewActivity(curSession.Activity, curWindows)
		if len(newActivity) > 0 {
			err := db.UpdateNewActivity(sbClient, curSession, newActivity)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = db.IncrementAll(sbClient, curSession, tracker)
		if err != nil {
			log.Fatal(err)
		}

		tracker.Reset()
	}
}
