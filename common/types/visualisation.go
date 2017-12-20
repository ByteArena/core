package types

import (
	"github.com/bytearena/core/common/utils/vector"
)

type VizMessage struct {
	GameID        string
	Objects       []VizMessageObject
	DebugPoints   [][2]float64
	DebugSegments [][2][2]float64
	Events        []VizMessageEvent
}

type VizMessageObject struct {
	Id          string
	Type        string
	Position    vector.Vector2
	Velocity    vector.Vector2
	Radius      float64
	Orientation float64

	PlayerInfo *PlayerInfo
}

type PlayerInfo struct {
	IsAlive    bool
	PlayerId   string
	PlayerName string
	Score      VizMessagePlayerScore
}

type VizMessagePlayerScore struct {
	Value int
}

type VizMessageEvent struct {
	Subject string
	Payload interface{}
}
