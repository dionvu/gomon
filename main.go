package main

import (
	"log"
	"time"

	"github.com/dionvu/gomon/db"
	"github.com/dionvu/gomon/session"
	sb "github.com/nedpals/supabase-go"
)

const (
	INCREMENT_INTERVAL = time.Second * 2
	KEYBOARD_FILE      = "/dev/input/event19"
	TRACKPAD_FILE      = "/dev/input/event11"
	MOUSE_FILE         = "/dev/input/event21"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)

	sbClient := sb.CreateClient(db.LoadSecret())

	tracker := session.NewTracker(KEYBOARD_FILE, MOUSE_FILE, TRACKPAD_FILE)
	defer tracker.Close()

	tracker.ListenAll()

	for {
		var curSession session.Session

		curSession, err := db.GetCurrentSession(sbClient)
		if err != nil {
			_, err = db.AddSession(sbClient)
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

		err = db.IncrementAll(sbClient, curSession, tracker)
		if err != nil {
			log.Fatal(err)
		}

		tracker.Reset()
	}
}
