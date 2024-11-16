package input

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
)

type Event struct {
	TimeSec  uint32
	TimeUsec uint32
	Type     EventType
	Code     EventCode
	Value    EventValue
}

type EventType uint16

const (
	EV_KEY EventType = 0x01
	EV_REL EventType = 0x02
)

type EventCode uint16

const (
	BTN_LEFT   EventCode = 0x110
	BTN_RIGHT  EventCode = 0x111
	BTN_MIDDLE EventCode = 0x112
	REL_X      EventCode = 0x00
	REL_Y      EventCode = 0x01
)

type EventValue int32

const (
	KEY_HOLD EventValue = 2
	KEY_DOWN EventValue = 1
	KEY_LIFT EventValue = 0
)

// The size that an event takes up in bytes.
const SIZE = 24

// Parses an input event from bytes
// revieved from /dev/input/eventX.
func From(b []byte) (Event, error) {
	event := Event{}

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

func (event Event) IsRightClick() bool {
	return event.Type == EV_KEY &&
		event.Code == BTN_RIGHT &&
		event.Value == KEY_DOWN
}

func (event Event) IsLeftClick() bool {
	return event.Type.Equals(EV_KEY) &&
		event.Code.Equals(BTN_LEFT) &&
		event.Value == KEY_DOWN
}

func (event Event) IsMiddleClick() bool {
	return event.Type.Equals(EV_KEY) &&
		event.Code.Equals(BTN_MIDDLE) &&
		event.Value == KEY_DOWN
}

func (event Event) IsMouseMove() bool {
	return event.Type.Equals(EV_REL) &&
		event.Code.Equals(REL_X, REL_Y)
}

func (event Event) IsKeyboardPress() bool {
	validKey := 0 <= event.Code && event.Code <= 248

	return event.Type.Equals(EV_KEY) && validKey && event.Value == KEY_LIFT
}

func (value EventValue) Meter() float64 {
	const conversion = 0.0000244
	return float64(value) * conversion
}

func (value EventValue) Abs() EventValue {
	return EventValue(math.Abs(float64(value)))
}

func (code EventCode) IsRelX() bool {
	return code == REL_X
}

func (code EventCode) IsRelY() bool {
	return code == REL_Y
}
