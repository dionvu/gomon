package archive

import (
	"time"

	"github.com/dionvu/gomon/session"
	"github.com/google/uuid"
)

type Archive struct {
	Id       string            `json:"id"`
	Created  time.Time         `json:"created"`
	Date     string            `json:"date"`
	Sessions []session.Session `json:"sessions"`
}

func NewArchive(sessions []session.Session) Archive {
	return Archive{
		Id:       uuid.NewString(),
		Created:  time.Now(),
		Date:     time.Now().Format(time.DateOnly),
		Sessions: sessions,
	}
}
