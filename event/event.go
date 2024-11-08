package event

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// https://www.kernel.org/doc/html/v4.17/input/event-codes.html
// /nix/store/hnxrl4wbjkpfkfdq6p01qgxhlg8lhzya-linux-headers-static-6.5/include/linux/input-event-codes.h

type InputEvent struct {
	TimeSec  uint32
	TimeUsec uint32
	Type     EventType
	Code     EventCode
	Value    int32
}

type EventType uint16

type EventCode uint16

const (
	EV_KEY EventType = 0x01
	EV_REL EventType = 0x02
)

const (
	BTN_LEFT  EventCode = 0x110
	BTN_RIGHT EventCode = 0x111
)

const SIZE = 24

func From(b []byte) (InputEvent, error) {
	event := InputEvent{}

	if len(b) < 20 {
		return event, errors.New("Events are not shorter than 20 bytes")
	}

	binary.Read(bytes.NewReader(b[0:8]), binary.LittleEndian, &event.TimeSec)
	binary.Read(bytes.NewReader(b[8:16]), binary.LittleEndian, &event.TimeUsec)
	binary.Read(bytes.NewReader(b[16:18]), binary.LittleEndian, &event.Type)
	binary.Read(bytes.NewReader(b[18:20]), binary.LittleEndian, &event.Code)
	binary.Read(bytes.NewReader(b[20:]), binary.LittleEndian, &event.Value)

	return event, nil
}

func (event InputEvent) IsRightClick() bool {
	return event.Type == EV_KEY &&
		event.Code == BTN_RIGHT &&
		event.Value == 1
}

func (event InputEvent) IsLeftClick() bool {
	return event.Type == EV_KEY &&
		event.Code == BTN_LEFT &&
		event.Value == 1
}

// Returns if code is equal to any
// of the passed EventCodes.
func (code EventCode) Equals(codes ...EventCode) bool {
	for _, c := range codes {
		if code == c {
			return true
		}
	}
	return false
}

// Returns if typ is equal to any
// of the passed inputTypes.
func (typ EventType) Equals(types ...EventType) bool {
	for _, t := range types {
		if typ == t {
			return true
		}
	}
	return false
}
