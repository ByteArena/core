package deathmatch

import (
	"github.com/go-gl/mathgl/mgl64"

	"github.com/bytearena/box2d"

	"github.com/bytearena/core/common/utils/vector"
)

type PhysicalBody struct {
	body *box2d.B2Body

	// 2 dimensional transform for points in space
	pointTransformIn  *mgl64.Mat4
	pointTransformOut *mgl64.Mat4

	// 1 dimensional transform for distance; could be infered from 2 dimensional transforms, but easier / faster that way
	distanceScaleIn  float64
	distanceScaleOut float64

	// 1 dimensional transform for time
	timeScaleIn  float64
	timeScaleOut float64

	maxSpeed           float64 // expressed in m/tick (agent referential)
	maxAngularVelocity float64 // expressed in rad/tick (agent referential)
	dragForce          float64 // expressed in m/tick (agent referential)
	static             bool
	skipThisTurn       bool
}

func (p *PhysicalBody) GetBody() *box2d.B2Body {
	return p.body
}

func (p *PhysicalBody) SetBody(body *box2d.B2Body) *PhysicalBody {
	p.body = body
	return p
}

func (p PhysicalBody) GetPhysicalReferentialPosition() vector.Vector2 {
	v := p.body.GetPosition()
	return vector.MakeVector2(v.X, v.Y)
}

func (p PhysicalBody) GetPosition() vector.Vector2 {
	v := p.body.GetPosition()
	return vector.MakeVector2(v.X, v.Y).Transform(p.pointTransformOut)
}

func (p *PhysicalBody) SetPosition(v vector.Vector2) *PhysicalBody {
	p.body.SetTransform(v.Transform(p.pointTransformIn).ToB2Vec2(), p.GetOrientation())
	return p
}

func (p *PhysicalBody) SetPositionInPhysicalScale(v vector.Vector2) *PhysicalBody {
	p.body.SetTransform(v.ToB2Vec2(), p.GetOrientation())
	return p
}

func (p PhysicalBody) GetPhysicalReferentialVelocity() vector.Vector2 {
	v := p.body.GetLinearVelocity()
	return vector.MakeVector2(v.X, v.Y)
}

func (p PhysicalBody) GetVelocity() vector.Vector2 {

	v := p.body.GetLinearVelocity()

	// In Box2D, velocity is expressed in m/s in physics scale
	// In Game, velocity is expressed in m/tick in agent scale

	// Box2D => game : v * timescaleOut * transformOut

	return vector.
		MakeVector2(v.X, v.Y).
		Scale(p.timeScaleOut).
		Transform(p.pointTransformOut)
}

func (p *PhysicalBody) SetVelocity(v vector.Vector2) *PhysicalBody {

	// In Box2D, velocity is expressed in m/s in physics scale
	// In Game, velocity is expressed in m/tick in agent scale

	// Game => Box2D : v * timeScaleIn * transformIn

	box2dvelocity := v.
		Scale(p.timeScaleIn).
		Transform(p.pointTransformIn).
		ToB2Vec2()

	p.body.SetLinearVelocity(box2dvelocity)

	return p
}

func (p PhysicalBody) GetPhysicalReferentialOrientation() float64 {
	return p.body.GetAngle()
}

func (p PhysicalBody) GetOrientation() float64 {
	return p.body.GetAngle() // no transform on angles
}

func (p *PhysicalBody) SetOrientation(angle float64) *PhysicalBody {
	// Could also be implemented using torque; see http://www.iforce2d.net/b2dtut/rotate-to-angle
	p.body.SetTransform(p.body.GetPosition(), angle)
	return p
}

func (p PhysicalBody) GetPhysicalReferentialRadius() float64 {
	return p.body.GetFixtureList().GetShape().GetRadius()
}

func (p PhysicalBody) GetRadius() float64 {
	// here we suppose that the body is always a circle
	return p.body.GetFixtureList().GetShape().GetRadius() * p.distanceScaleOut
}

func (p PhysicalBody) GetMaxSpeed() float64 {
	return p.maxSpeed
}

func (p *PhysicalBody) SetMaxSpeed(maxSpeed float64) *PhysicalBody {
	p.maxSpeed = maxSpeed
	return p
}

func (p PhysicalBody) GetMaxAngularVelocity() float64 {
	return p.maxAngularVelocity
}

func (p *PhysicalBody) SetMaxAngularVelocity(maxAngularVelocity float64) *PhysicalBody {
	p.maxAngularVelocity = maxAngularVelocity
	return p
}

func (p PhysicalBody) GetDragForce() float64 {
	return p.dragForce
}

func (p *PhysicalBody) SetDragForce(dragForce float64) *PhysicalBody {
	p.dragForce = dragForce
	return p
}
