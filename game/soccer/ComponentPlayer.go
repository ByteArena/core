package soccer

import "github.com/bytearena/core/common/types"

type Player struct {
	// Populated by systemScore
	Score int

	Agent *types.Agent
}
