package db

import (
	"errors"
	"log"
	"os"
	"slices"
	"time"

	"github.com/dionvu/gomon/archive"
	"github.com/dionvu/gomon/hypr"
	"github.com/dionvu/gomon/session"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	sb "github.com/nedpals/supabase-go"
)

const (
	TABLE_SESSIONS    = "sessions"
	TABLE_ARCHIVE     = "archives"
	COL_ACTIVITY      = "activity"
	COL_ID            = "id"
	COL_KEYPRESSES    = "key_presses"
	COL_LEFT_CLICKS   = "left_clicks"
	COL_RIGHT_CLICKS  = "right_clicks"
	COL_MIDDLE_CLICKS = "middle_clicks"
	COL_MOUSE_X       = "x_mouse_movement"
	COL_MOUSE_Y       = "y_mouse_movement"
	COL_SESSIONS      = "sessions"
	COL_START         = "start"
	COL_END           = "end"
)

const (
	SESSION_INTERVALS  = time.Minute * 30
	INCREMENT_INTERVAL = time.Minute
)

// Adds a session with given activity, the sessions start and end is marked
// by the current time rounded down and up an hour, respectively.
func AddSession(client *sb.Client, activity []session.Activity) (interface{}, error) {
	ses := session.Session{
		Id:       uuid.NewString(),
		Start:    time.Now().Truncate(SESSION_INTERVALS),
		End:      time.Now().Truncate(SESSION_INTERVALS).Add(SESSION_INTERVALS),
		Activity: activity,
	}

	var res []interface{}

	_, err := GetCurrentSession(client)
	if err == nil {
		return nil, errors.New("There is an existing session in the database")
	}

	err = client.DB.From(TABLE_SESSIONS).Insert(ses).Execute(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

// Gets the session that with intervals start
// and end that match the current time.
func GetCurrentSession(client *sb.Client) (session.Session, error) {
	var res []session.Session
	err := client.DB.From(TABLE_SESSIONS).Select("*").Eq("start", FormatTime(time.Now().Truncate(SESSION_INTERVALS))).Execute(&res)
	if err != nil {
		return session.Session{}, err
	}

	if len(res) == 0 {
		return session.Session{}, errors.New("No current session in database")
	}

	return res[0], nil
}

func GetAllSessionsToday(client *sb.Client) ([]session.Session, error) {
	var res []session.Session

	startOfDay := FormatTime(time.Now().Truncate(24 * time.Hour))
	endOfDay := FormatTime(time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour))

	err := client.DB.From(TABLE_SESSIONS).Select("*").Gte(COL_START, startOfDay).Lt(COL_END, endOfDay).Execute(&res)
	if err != nil {
		return res, err
	}

	slices.SortFunc(res, func(a, b session.Session) int {
		return a.Start.Hour() - b.Start.Hour()
	})

	return res, nil
}

func DropSessions(client *sb.Client, before time.Time) (interface{}, error) {
	var res interface{}

	err := client.DB.From(TABLE_SESSIONS).Delete().Lt(COL_END, FormatTime(before)).Execute(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func ArchivePastSessions(client *sb.Client) (interface{}, error) {
	var res []interface{}
	exists := map[string]bool{}

	arc, err := GetCurrentArchive(client)
	if err != nil {
		AddNewArchive(client)
		arc, err = GetCurrentArchive(client)
		if err != nil {
			return res, err
		}
	}

	for _, ses := range arc.Sessions {
		exists[ses.Id] = true
	}

	sessions, err := GetAllSessionsToday(client)
	if err != nil {
		return res, err
	}

	for _, ses := range sessions {
		// fmt.Println(ses.End)

		if ses.End.Before(time.Now()) && !exists[ses.Id] {
			ArchiveSession(client, ses)
		}
	}

	return res, nil
}

func ArchiveSession(client *sb.Client, ses session.Session) (interface{}, error) {
	var res []interface{}

	arc, err := GetCurrentArchive(client)
	if err != nil {
		AddNewArchive(client)
		arc, err = GetCurrentArchive(client)
		if err != nil {
			return res, err
		}
	}

	arc.Sessions = append(arc.Sessions, ses)

	updatedData := map[string]interface{}{
		COL_SESSIONS: arc.Sessions,
	}

	err = client.DB.From(TABLE_ARCHIVE).Update(updatedData).Eq(COL_ID, arc.Id).Execute(&res)

	return res, nil
}

func AddNewArchive(client *sb.Client) (interface{}, error) {
	var res interface{}

	arc := archive.NewArchive([]session.Session{})

	err := client.DB.From(TABLE_ARCHIVE).Insert(arc).Execute(&res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func GetCurrentArchive(client *sb.Client) (archive.Archive, error) {
	var res []archive.Archive

	err := client.DB.From(TABLE_ARCHIVE).Select("*").Eq("date", time.Now().Format(time.DateOnly)).Execute(&res)
	if err != nil {
		return archive.Archive{}, err
	}

	if len(res) == 0 {
		return archive.Archive{}, errors.New("No current archive in database")
	}

	return res[0], nil
}

// Given the current session's id, increments the time spent for each activity only
// if the same activity is found in windows (the user's current windows).
func IncrementActivityTime(client *sb.Client, currSession session.Session, activity []session.Activity,
	windows []hypr.Window, add time.Duration,
) (interface{}, error) {
	if currSession.Id == "" {
		return nil, nil
	}

	var res interface{}
	exists := map[hypr.Window]bool{}

	for _, win := range windows {
		exists[win] = true
	}

	for i, a := range activity {
		if exists[a.Window] {
			activity[i].TimeSpentMin += add.Minutes()
		}
	}

	updatedData := map[string]interface{}{
		COL_ACTIVITY: activity,
	}

	err := client.DB.From(TABLE_SESSIONS).Update(updatedData).Eq("id", currSession.Id).Execute(res)
	if err != nil {
		return nil, err
	}

	return res, err
}

// Updates the passed current session with the passed new activity.
func UpdateNewActivity(client *sb.Client, currSession session.Session, newActivity []session.Activity) error {
	var res interface{}

	for _, activity := range newActivity {
		currSession.Activity = append(currSession.Activity, activity)
	}

	updatedActivity := map[string]interface{}{
		COL_ACTIVITY: currSession.Activity,
	}

	err := client.DB.From(TABLE_SESSIONS).Update(updatedActivity).Eq(COL_ID, currSession.Id).Execute(res)
	if err != nil {
		return err
	}

	return nil
}

func IncrementClickCount(client *sb.Client, curSession session.Session, tracker session.Tracker) error {
	var res interface{}

	updatedStats := map[string]interface{}{
		COL_LEFT_CLICKS:   curSession.LeftClicks + tracker.LeftClicks,
		COL_RIGHT_CLICKS:  curSession.RightClicks + tracker.RightCLicks,
		COL_MIDDLE_CLICKS: curSession.MiddleClicks + tracker.MiddleClicks,
	}

	err := client.DB.From(TABLE_SESSIONS).Update(updatedStats).Eq(COL_ID, curSession.Id).Execute(&res)
	if err != nil {
		return err
	}

	return nil
}

func IncrementKeyboardPressCount(client *sb.Client, curSession session.Session, tracker session.Tracker) error {
	var res interface{}

	updatedStats := map[string]interface{}{
		COL_KEYPRESSES: curSession.KeyPresses + tracker.KeyboardPresses,
	}

	err := client.DB.From(TABLE_SESSIONS).Update(updatedStats).Eq(COL_ID, curSession.Id).Execute(&res)
	if err != nil {
		return err
	}

	return nil
}

func IncrementMouseMovement(client *sb.Client, curSession session.Session, tracker session.Tracker) error {
	var res interface{}

	updatedStats := map[string]interface{}{
		COL_MOUSE_X: curSession.XMovement + tracker.XMovement,
		COL_MOUSE_Y: curSession.YMovement + tracker.YMovement,
	}

	err := client.DB.From(TABLE_SESSIONS).Update(updatedStats).Eq(COL_ID, curSession.Id).Execute(&res)
	if err != nil {
		return err
	}

	return nil
}

func IncrementAll(client *sb.Client, curSession session.Session, tracker session.Tracker) error {
	err := IncrementClickCount(client, curSession, tracker)
	if err != nil {
		return err
	}

	err = IncrementKeyboardPressCount(client, curSession, tracker)
	if err != nil {
		return err
	}

	err = IncrementMouseMovement(client, curSession, tracker)
	if err != nil {
		return err
	}

	return nil
}

// Loads database secrets, returning the
// db url and the api key, else panics!
func LoadSecret() (string, string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	return os.Getenv("DB_URL"), os.Getenv("API_KEY")
}

// "2006-01-02T15:04:05Z07:00"
func FormatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
