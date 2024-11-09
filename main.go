package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dionvu/temp/db"
	"github.com/dionvu/temp/hypr"
	"github.com/dionvu/temp/session"
	"github.com/joho/godotenv"
	sb "github.com/nedpals/supabase-go"
)

var (
	dbUrl  string
	apiKey string
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbUrl = os.Getenv("DB_URL")
	apiKey = os.Getenv("API_KEY")

	client := sb.CreateClient(dbUrl, apiKey)

	if curSession, err := db.GetCurrentSession(client); err == nil {
		for {
			curSession, err = db.GetCurrentSession(client)
			if err != nil {
				log.Fatal(err)
			}

			time.Sleep(time.Minute)

			windows, err := hypr.Windows()
			if err != nil {
				log.Fatal(err)
			}

			_, err = db.IncrementActivityTime(client, curSession.Id, curSession.Activity, windows, time.Minute)
			if err != nil {
				log.Fatal(err)
			}

			newActivity := session.FilterNewActivity(curSession.Activity, windows)
			if len(newActivity) > 0 {
				db.UpdateNewActivity(client, curSession, newActivity)
			}

			fmt.Println(newActivity)
		}
	}

	windows, err := hypr.Windows()
	if err != nil {
		log.Fatal(err)
	}

	activity := session.NewActivity(windows)

	_, err = db.AddHourSession(client, activity)
	if err != nil {
		log.Fatal(err)
	}

	curSession, err := db.GetCurrentSession(client)
	if err != nil {
		log.Fatal(err)
	}

	for {
		curSession, err = db.GetCurrentSession(client)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(time.Minute)

		windows, err := hypr.Windows()
		if err != nil {
			log.Fatal(err)
		}
		_, err = db.IncrementActivityTime(client, curSession.Id, curSession.Activity, windows, time.Minute)
		if err != nil {
			log.Fatal(err)
		}

		newActivity := session.FilterNewActivity(curSession.Activity, windows)
		if len(newActivity) > 0 {
			db.UpdateNewActivity(client, curSession, newActivity)
		}

		fmt.Println(newActivity)
	}
}
