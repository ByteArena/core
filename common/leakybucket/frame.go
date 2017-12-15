package leakybucket

import "time"

type Framenum int

type Frame struct {
	num     Framenum
	time    time.Time
	payload string
}

func makeFrame(num Framenum, payload string) Frame {
	return Frame{
		num:     num,
		time:    time.Now(),
		payload: payload,
	}
}

func (frame Frame) GetNum() Framenum {
	return frame.num
}

func (frame Frame) GetTime() time.Time {
	return frame.time
}

func (frame Frame) GetPayload() string {
	return frame.payload
}
