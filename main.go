package main

import (
	"log"
	"os"
	"time"

	"github.com/dionvu/temp/db"
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

	// activity, err := session.NewActivity()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = db.AddHourSession(client, activity)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	s, err := db.GetCurrentSession(client)
	if err != nil {
		log.Fatal(err)
	}

	db.IncrementActivityTime(client, s.Id, s.Activity, time.Minute)
}
