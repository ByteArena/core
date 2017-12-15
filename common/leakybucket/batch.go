package leakybucket

import "time"

type Batchnum int

type Batch struct {
	num       Batchnum
	fromFrame Framenum
	toFrame   Framenum
	fromTime  time.Time
	toTime    time.Time
	frames    []Frame
	numframes int
}

func newBatch(num Batchnum, framesPerBatch int) *Batch {
	return &Batch{
		num:       num,
		fromFrame: -1,
		toFrame:   -1,
		frames:    make([]Frame, framesPerBatch),
		numframes: 0,
	}
}

func (batch *Batch) size() int {
	return batch.numframes
}

func (batch *Batch) addFrame(frame Frame) {
	if batch.fromFrame == -1 {
		batch.fromFrame = frame.GetNum()
		batch.fromTime = frame.GetTime()
	}

	batch.toFrame = frame.GetNum()
	batch.toTime = frame.GetTime()

	batch.frames[batch.numframes] = frame
	batch.numframes++
}

func (batch *Batch) GetFrames() []Frame {
	return batch.frames
}
