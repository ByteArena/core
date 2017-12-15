package leakybucket

import (
	"time"
)

// Responsible for batching and buffering
// bundle frames in batches and store them in buffer

type Bucket struct {
	startTime        time.Time
	framesPerBatch   int
	batchesToKeep    int
	batches          *deque
	framenum         Framenum
	batchnum         Batchnum
	batchInTheMaking *Batch
	onBatch          func(batch Batch, bucket *Bucket)
}

func NewBucket(framesPerBatch int, batchesToKeep int, onBatch func(batch Batch, bucket *Bucket)) *Bucket {
	return &Bucket{
		framesPerBatch:   framesPerBatch,
		batchesToKeep:    batchesToKeep,
		batches:          newDeque(),
		framenum:         0,
		batchnum:         0,
		batchInTheMaking: newBatch(0, framesPerBatch),
		onBatch:          onBatch,
	}
}

func (bucket *Bucket) AddFrame(payload string) {
	bucket.batchInTheMaking.addFrame(makeFrame(bucket.framenum, payload))
	bucket.framenum++

	if bucket.batchInTheMaking.size() >= bucket.framesPerBatch {
		newbatch := *bucket.batchInTheMaking // copy content
		bucket.batchInTheMaking = newBatch(bucket.batchnum, bucket.framesPerBatch)
		bucket.batchnum++
		bucket.batches.append(newbatch)

		// trimming if needed
		for bucket.batches.size() > bucket.batchesToKeep {
			bucket.batches.shift()
		}

		bucket.onBatch(newbatch, bucket)
	}
}

func (bucket *Bucket) GetBatches() []Batch {
	res := make([]Batch, bucket.batches.size())
	i := 0
	for _, batch := range bucket.batches.iterate() {
		res[i] = batch
		i++
	}

	return res
}
