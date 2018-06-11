package mailboxmessages

import "github.com/bytearena/ecs"

type YouHaveExitedTheMaze struct {
	Entity ecs.EntityID
}

func (msg YouHaveExitedTheMaze) Subject() string {
	return "exitedmaze"
}
