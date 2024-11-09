
// package main
//
// import (
// 	"fmt"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"time"
//
// 	"github.com/dionvu/temp/event"
// )
//
// const (
// 	EVENT_DIR        = "/dev/input"
// 	MOUSE_EVENT_FILE = "event19"
// )
//
// var (
// 	leftClicks  = 0
// 	rightClicks = 0
// )
//
// func main() {
// 	mouseEventFile, err := os.Open(filepath.Join(EVENT_DIR, MOUSE_EVENT_FILE))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer mouseEventFile.Close()
//
// 	b := make([]byte, event.SIZE)
//
// 	go func() {
// 		printCountsLoop()
// 	}()
//
// 	for {
// 		mouseEventFile.Read(b)
//
// 		event, err := event.From(b)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
//
// 		if event.IsLeftClick() {
// 			leftClicks += 1
// 		}
//
// 		if event.IsRightClick() {
// 			rightClicks += 1
// 		}
// 	}
// }
//
// func prinCounts() {
// 	fmt.Println("LC: ", leftClicks)
// 	fmt.Println("RC: ", rightClicks)
// }
//
// func printCountsLoop() {
// 	prinCounts()
//
// 	time.Sleep(time.Second * 5)
//
// 	printCountsLoop()
// }
