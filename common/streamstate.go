package common

import (
	"encoding/json"

	"github.com/bytearena/core/arenaserver"
	"github.com/bytearena/core/common/leakybucket"
	"github.com/bytearena/core/common/mq"
)

func StreamState(srv *arenaserver.Server, brokerclient mq.ClientInterface, arenaServerUUID string) {

	buk := leakybucket.NewBucket(
		srv.GetTicksPerSecond(),
		5, // keep 5 seconds of stream in buffer
		func(batch leakybucket.Batch, bucket *leakybucket.Bucket) {
			frames := batch.GetFrames()
			jsonbatch := make([]json.RawMessage, len(frames))
			for i, frame := range frames {
				jsonbatch[i] = json.RawMessage(frame.GetPayload())
			}

			brokerclient.Publish("viz", "message", jsonbatch)
		},
	)

	stateobserver := srv.SubscribeStateObservation()
	for {
		select {
		case <-stateobserver:
			{
				buk.AddFrame(string(srv.GetGame().GetVizFrameJson()))
			}
		}
	}
}
