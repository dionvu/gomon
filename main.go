package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/dionvu/temp/event"
)

const (
	EVENT_DIR        = "/dev/input"
	MOUSE_EVENT_FILE = "event19"
)

func main() {
	mouseEventFile, err := os.Open(filepath.Join(EVENT_DIR, MOUSE_EVENT_FILE))
	if err != nil {
		log.Fatal(err)
	}
	defer mouseEventFile.Close()

	b := make([]byte, event.SIZE)

	for {
		mouseEventFile.Read(b)

		event, err := event.From(b)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(event.IsRightClick())
	}
}
