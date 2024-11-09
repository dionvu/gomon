package db

import (
	"errors"
	"time"

	"github.com/dionvu/temp/hypr"
	"github.com/dionvu/temp/session"
	"github.com/google/uuid"
	sb "github.com/nedpals/supabase-go"
)

const (
	TABLE_SESSIONS  = "sessions"
	COLUMN_ACTIVITY = "activity"
	COLUMN_ID       = "id"
)

// Adds a session with given activity, the sessions start and end is marked
// by the current time rounded down and up an hour, respectively.
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

// Gets the session that with intervals start
// and end that match the current time.
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

// Given the current session's id, increments the time spent for each activity only
// if the same activity is found in windows (the user's current windows).
func IncrementActivityTime(client *sb.Client, id string, activity []session.Activity, windows []hypr.Window, add time.Duration) (interface{}, error) {
	var res interface{}
	m := map[string]bool{}

	for _, win := range windows {
		m[win.Id] = true
	}

	for i, a := range activity {
		if m[a.Window.Id] {
			activity[i].TimeSpentMin += add.Minutes()
		}
	}

	updatedData := map[string]interface{}{
		COLUMN_ACTIVITY: activity,
	}

	err := client.DB.From(TABLE_SESSIONS).Update(updatedData).Eq("id", id).Execute(res)
	if err != nil {
		return nil, err
	}

	return res, err
}

func UpdateNewActivity(client *sb.Client, currSession session.Session, newActivity []session.Activity) error {
	var res interface{}

	for _, activity := range newActivity {
		currSession.Activity = append(currSession.Activity, activity)
	}

	updatedActivity := map[string]interface{}{
		COLUMN_ACTIVITY: currSession.Activity,
	}

	err := client.DB.From(TABLE_SESSIONS).Update(updatedActivity).Eq(COLUMN_ID, currSession.Id).Execute(res)
	if err != nil {
		return err
	}

	return nil
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
