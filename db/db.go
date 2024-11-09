package db

import (
	"errors"
	"time"

	"github.com/dionvu/temp/session"
	"github.com/google/uuid"
	sb "github.com/nedpals/supabase-go"
)

const TABLE_SESSIONS = "sessions"

func AddHourSession(client *sb.Client, activity []session.Activity) (interface{}, error) {
	hour := session.Session{
		Id:       uuid.New().String(),
		Start:    time.Now().Truncate(time.Hour),
		End:      time.Now().Truncate(time.Hour).Add(time.Hour),
		Activity: activity,
	}

	var res []interface{}

	_, err := GetCurrentSession(client)
	if err == nil {
		return nil, errors.New("There is an existing session in the database")
	}

	err = client.DB.From(TABLE_SESSIONS).Insert(hour).Execute(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func GetCurrentSession(client *sb.Client) (session.Session, error) {
	var res []session.Session
	err := client.DB.From(TABLE_SESSIONS).Select("*").Eq("start", formatTime(time.Now().Truncate(time.Hour))).Execute(&res)
	if err != nil {
		return session.Session{}, err
	}

	if len(res) == 0 {
		return session.Session{}, errors.New("No current session in database")
	}

	return res[0], nil
}

func IncrementActivityTime(client *sb.Client, id string, activity []session.Activity, add time.Duration) (interface{}, error) {
	for i := range activity {
		activity[i].TimeSpentMin += int(add.Minutes())
	}

	updatedData := struct {
		Activity []session.Activity `json:"activity"`
	}{
		Activity: activity,
	}

	var res interface{}
	err := client.DB.From("sessions").Update(updatedData).Eq("id", id).Execute(res)
	if err != nil {
		return nil, err
	}

	return res, err
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
