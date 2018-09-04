package soccer

import (
	"github.com/bytearena/core/common/utils/vector"
	"github.com/bytearena/ecs"
)

type agentPerception struct {
	Score    int                         `json:"score"`
	Velocity vector.Vector2              `json:"velocity"` // vecteur de force (direction, magnitude)
	Vision   []agentPerceptionVisionItem `json:"vision"`
}

var agentPerceptionVisionItemTag = struct {
	Agent    string
	Obstacle string
	Ball     string
}{
	Agent:    "agent",
	Obstacle: "obstacle",
	Ball:     "ball",
}

type agentPerceptionVisionItem struct {
	Tag      string         `json:"tag"`
	Center   vector.Point2  `json:"center"`
	Velocity vector.Vector2 `json:"velocity"`
	Radius   float64        `json:"radius"`
	EntityID ecs.EntityID   `json:"-"`
}

func newEmptyAgentPerception() *agentPerception {
	return &agentPerception{
		Vision: make([]agentPerceptionVisionItem, 0),
	}
}
