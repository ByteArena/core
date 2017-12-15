package deathmatch

import (
	"github.com/bytearena/ecs"

	"github.com/bytearena/core/common/utils"
	"github.com/bytearena/core/common/utils/vector"
)

var CollisionGroup = struct {
	Ground     utils.Tag
	Obstacle   utils.Tag
	Agent      utils.Tag
	Projectile utils.Tag
}{
	Ground:     utils.MakeTag(0),
	Obstacle:   utils.MakeTag(1),
	Agent:      utils.MakeTag(2),
	Projectile: utils.MakeTag(3),
}

type collisionScriptFunc func(game *DeathmatchGame, entityID ecs.EntityID, otherEntityID ecs.EntityID, collidableAspect *Collidable, otherCollidableAspectB *Collidable, point vector.Vector2)

type Collidable struct {
	collisiongroup      utils.Tag
	collideswith        utils.Tag
	collisionScriptFunc collisionScriptFunc
}

func NewCollidable(collisiongroup utils.Tag, collideswith utils.Tag) *Collidable {
	return &Collidable{
		collisiongroup: collisiongroup,
		collideswith:   collideswith,
	}
}

func (collidable *Collidable) SetCollisionScriptFunc(f collisionScriptFunc) *Collidable {
	collidable.collisionScriptFunc = f
	return collidable
}

func (collidable *Collidable) MayCollideWith(othercollidable *Collidable) bool {
	return collidable.collideswith.Includes(othercollidable.collisiongroup)
}

func (collidable *Collidable) CollisionScript(game *DeathmatchGame, entityID ecs.EntityID, otherEntityID ecs.EntityID, collidableAspect *Collidable, otherCollidableAspectB *Collidable, point vector.Vector2) {
	if collidable.collisionScriptFunc == nil {
		return
	}

	collidable.collisionScriptFunc(game, entityID, otherEntityID, collidableAspect, otherCollidableAspectB, point)
}
