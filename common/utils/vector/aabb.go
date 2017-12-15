package vector

import (
	"github.com/bytearena/box2d"
	"github.com/go-gl/mathgl/mgl64"
)

type AABB struct {
	LowerBound Vector2
	UpperBound Vector2
}

func GetAABBForPointList(points ...Vector2) AABB {

	var minX = 10000000000.0
	var minY = 10000000000.0
	var maxX = -10000000000.0
	var maxY = -10000000000.0

	for _, point := range points {
		x, y := point.Get()
		if x < minX {
			minX = x
		}

		if y < minY {
			minY = y
		}

		if x > maxX {
			maxX = x
		}

		if y > maxY {
			maxY = y
		}
	}

	return AABB{
		LowerBound: MakeVector2(minX, minY),
		UpperBound: MakeVector2(maxX, maxY),
	}
}

func (aabb AABB) Overlaps(otheraabb AABB) bool {
	ourMinX, ourMinY := aabb.LowerBound.Get()
	ourMaxX, ourMaxY := aabb.UpperBound.Get()
	otherMinX, otherMinY := otheraabb.LowerBound.Get()
	otherMaxX, otherMaxY := otheraabb.UpperBound.Get()

	return !(ourMinX > otherMaxX ||
		ourMinY > otherMaxY ||
		ourMaxX < otherMinX ||
		ourMaxY < otherMinY)
}

func (aabb AABB) Transform(transform *mgl64.Mat4) AABB {
	aabb.LowerBound = aabb.LowerBound.Transform(transform)
	aabb.UpperBound = aabb.UpperBound.Transform(transform)
	return aabb
}

func (aabb AABB) ToB2AABB() box2d.B2AABB {
	b2aabb := box2d.MakeB2AABB()
	b2aabb.LowerBound = aabb.LowerBound.ToB2Vec2()
	b2aabb.UpperBound = aabb.UpperBound.ToB2Vec2()

	return b2aabb
}
