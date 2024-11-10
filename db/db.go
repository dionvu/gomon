package db

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/dionvu/gomon/archive"
	"github.com/dionvu/gomon/hypr"
	"github.com/dionvu/gomon/session"
	"github.com/dionvu/gomon/util"
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
	COL_MOUSE_X       = "mouse_movement_meter_x"
	COL_MOUSE_Y       = "mouse_movement_meter_y"
)

// Adds a session with given activity, the sessions start and end is marked
// by the current time rounded down and up an hour, respectively.
func AddHourSession(client *sb.Client, activity []session.Activity) (interface{}, error) {
	hour := session.Session{
		Id:       uuid.NewString(),
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

func ArchieveSession(client *sb.Client, ses session.Session) (interface{}, error) {
	var res []interface{}

	var arc archive.Archive

	arc, err := GetCurrentArchive(client)
	if err != nil {
		arc = archive.NewArchive([]session.Session{ses})

		err = client.DB.From(TABLE_ARCHIVE).Insert(arc).Execute(&res)
		if err != nil {
			return res, err
		}
	}

	arc.Sessions = append(arc.Sessions, ses)

	updatedData := map[string]interface{}{
		"sessions": arc.Sessions,
	}

	err = client.DB.From(TABLE_ARCHIVE).Update(updatedData).Eq(COL_ID, arc.Id).Execute(&res)

	return res, nil
}

// Gets the session that with intervals start
// and end that match the current time.
func GetCurrentSession(client *sb.Client) (session.Session, error) {
	var res []session.Session
	err := client.DB.From(TABLE_SESSIONS).Select("*").Eq("start", util.FormatTime(time.Now().Truncate(time.Hour))).Execute(&res)
	if err != nil {
		return session.Session{}, err
	}

	if len(res) == 0 {
		return session.Session{}, errors.New("No current session in database")
	}

	return res[0], nil
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

func IncrementClickCount(client *sb.Client, curSession session.Session, left uint, right uint, middle uint) error {
	var res interface{}

	updatedStats := map[string]interface{}{
		COL_LEFT_CLICKS:   curSession.LeftClicks + left,
		COL_RIGHT_CLICKS:  curSession.RightClicks + right,
		COL_MIDDLE_CLICKS: curSession.MiddleClicks + middle,
	}

	err := client.DB.From(TABLE_SESSIONS).Update(updatedStats).Eq(COL_ID, curSession.Id).Execute(&res)
	if err != nil {
		return err
	}

	return nil
}

func IncrementKeyboardPressCount(client *sb.Client, curSession session.Session, count uint) error {
	var res interface{}

	updatedStats := map[string]interface{}{
		COL_KEYPRESSES: curSession.KeyPresses + count,
	}

	err := client.DB.From(TABLE_SESSIONS).Update(updatedStats).Eq(COL_ID, curSession.Id).Execute(&res)
	if err != nil {
		return err
	}

	return nil
}

func IncrementMouseMovement(client *sb.Client, curSession session.Session, xMeter float64, yMeter float64) error {
	var res interface{}

	updatedStats := map[string]interface{}{
		COL_MOUSE_X: curSession.XMovementMeter + xMeter,
		COL_MOUSE_Y: curSession.YMovementMeter + yMeter,
	}

	err := client.DB.From(TABLE_SESSIONS).Update(updatedStats).Eq(COL_ID, curSession.Id).Execute(&res)
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
