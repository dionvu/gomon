package input

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"slices"
)

type (
	EventType  uint16
	EventCode  uint16
	EventValue int32
)

type Event struct {
	TimeSec  uint32
	TimeUsec uint32
	Type     EventType
	Code     EventCode
	Value    EventValue
}

const (
	EV_KEY EventType = 0x01
	EV_REL EventType = 0x02
	EV_ABS EventType = 0x03
)

const (
	BTN_LEFT           EventCode = 0x110
	BTN_RIGHT          EventCode = 0x111
	BTN_MIDDLE         EventCode = 0x112
	BTN_TOOL_FINGER    EventCode = 0x145
	BTN_TOOL_DOUBLETAP EventCode = 0x14d
	BTN_TOUCH          EventCode = 0x14a

	REL_X EventCode = 0x00
	REL_Y EventCode = 0x01

	ABS_X              EventCode = 0x00
	ABS_Y              EventCode = 0x01
	ABS_MT_TRACKING_ID EventCode = 0x39
)

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

func (code EventCode) Equals(codes ...EventCode) bool {
	return slices.Contains(codes, code)
}

func (typ EventType) Equals(types ...EventType) bool {
	return slices.Contains(types, typ)
}

func (val EventValue) Equals(values ...EventValue) bool {
	return slices.Contains(values, val)
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

func (event Event) IsTrackpadLeftClick() bool {
	return event.Type == EV_KEY &&
		event.Code == BTN_TOOL_FINGER &&
		event.Value == KEY_DOWN
}

func (event Event) IsTrackpadRightClick() bool {
	return event.Type == EV_KEY &&
		event.Code == BTN_TOOL_DOUBLETAP &&
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
	validKey := event.Code <= 248

	return event.Type.Equals(EV_KEY) && validKey && event.Value == KEY_LIFT
}

func (event Event) IsTrackpadMove() bool {
	return event.Type.Equals(EV_ABS) &&
		event.Code.Equals(ABS_X, ABS_Y)
}

func (event Event) IsTrackpadTouch() bool {
	return event.Type.Equals(EV_KEY) &&
		event.Code.Equals(BTN_TOUCH) &&
		event.Value.Equals(KEY_DOWN)
}

func (event Event) IsTrackpadLift() bool {
	return event.Type.Equals(EV_ABS) &&
		event.Code.Equals(ABS_MT_TRACKING_ID) &&
		event.Value.Equals(-1)
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
