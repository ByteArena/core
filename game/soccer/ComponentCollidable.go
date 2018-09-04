package soccer

import (
	"github.com/bytearena/ecs"

	"github.com/bytearena/core/common/utils"
	"github.com/bytearena/core/common/utils/vector"
)

var CollisionGroup = struct {
	Obstacle utils.Tag
	Agent    utils.Tag
	Sensor   utils.Tag
	Ball     utils.Tag
}{
	Obstacle: utils.MakeTag(0),
	Agent:    utils.MakeTag(1),
	Sensor:   utils.MakeTag(2),
	Ball:     utils.MakeTag(3),
}

type collisionScriptFunc func(game *SoccerGame, entityID ecs.EntityID, otherEntityID ecs.EntityID, collidableAspect *Collidable, otherCollidableAspectB *Collidable, point vector.Vector2)

type Collidable struct {
	collisiongroup      utils.Tag
	collideswith        utils.Tag
	collisionScriptFunc collisionScriptFunc
	isSensor            bool
}

func (collidable *Collidable) SetCollisionScriptFunc(f collisionScriptFunc) *Collidable {
	collidable.collisionScriptFunc = f
	return collidable
}

func (collidable *Collidable) MayCollideWith(othercollidable *Collidable) bool {
	return collidable.collideswith.Includes(othercollidable.collisiongroup)
}

func (collidable *Collidable) CollisionScript(game *SoccerGame, entityID ecs.EntityID, otherEntityID ecs.EntityID, collidableAspect *Collidable, otherCollidableAspectB *Collidable, point vector.Vector2) {
	if collidable.collisionScriptFunc == nil {
		return
	}

	collidable.collisionScriptFunc(game, entityID, otherEntityID, collidableAspect, otherCollidableAspectB, point)
}
