package trigo

import (
	"github.com/bytearena/core/common/utils/vector"
)

func collideOrientedRectangles(polyOne, polyTwo []vector.Vector2) []vector.Vector2 {

	points := make([]vector.Vector2, 0)

	if !DoClosedConvexPolygonsIntersect(polyOne, polyTwo) {
		return points
	}

	/*

	 C*--------------------*B
	  |                    |
	  |                    |
	  |                    |
	 D*--------------------*A

	*/

	for i := 0; i < 4; i++ {
		polyOneEdge := vector.MakeSegment2(polyOne[i], polyOne[(i+1)%4])
		for j := 0; j < 4; j++ {
			polyTwoEdge := vector.MakeSegment2(polyOne[j], polyOne[(j+1)%4])
			if collisionPoint, intersects, colinear, _ := SegmentSegmentIntersection(
				polyOneEdge,
				polyTwoEdge,
			); intersects {
				if !colinear {
					points = append(points, collisionPoint)
				} else {

					colinearIntersections := make([]vector.Vector2, 0)

					a1, a2 := polyOneEdge.Get()
					b1, b2 := polyTwoEdge.Get()
					if PointOnLineSegment(a1, b1, b2) {
						colinearIntersections = append(colinearIntersections, a1)
					}

					if PointOnLineSegment(a2, b1, b2) {
						colinearIntersections = append(colinearIntersections, a2)
					}

					if PointOnLineSegment(b1, a1, a2) {
						colinearIntersections = append(colinearIntersections, b1)
					}

					if PointOnLineSegment(b2, a1, a2) {
						colinearIntersections = append(colinearIntersections, b2)
					}

					if len(colinearIntersections) > 0 {
						centerOfMass, _ := ComputeCenterOfMass(colinearIntersections)
						points = append(points, centerOfMass)
					}
				}
			}
		}
	}

	if len(points) == 0 {
		// one is inside the other
		// compute the area and send the points of the smallest one
		polyOneHeight := vector.MakeSegment2(polyOne[0], polyOne[1])
		polyOneWidth := vector.MakeSegment2(polyOne[1], polyOne[2])

		polyTwoHeight := vector.MakeSegment2(polyTwo[0], polyTwo[1])
		polyTwoWidth := vector.MakeSegment2(polyTwo[1], polyTwo[2])

		polyOneAreaSq := polyOneWidth.LengthSq() * polyOneHeight.LengthSq()
		polyTwoAreaSq := polyTwoWidth.LengthSq() * polyTwoHeight.LengthSq()

		if polyTwoAreaSq < polyOneAreaSq {
			return polyTwo
		}

		return polyOne
	}

	return points
}
