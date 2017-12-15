package trigo

import (
	"github.com/bytearena/core/common/utils/vector"
)

func collideOrientedRectangleCircle(poly []vector.Vector2, center vector.Vector2, radius float64, trajectoryPointA, trajectoryPointB vector.Vector2, colliderRadius float64) []vector.Segment2 {

	collisionPositionSegments := make([]vector.Segment2, 0)

	points := make([]vector.Vector2, 0)

	// 2. On détermine les segments du polygone pour lesquels calculer une tangente

	trajectorySegment := vector.MakeSegment2(trajectoryPointA, trajectoryPointB)
	trajectoryLength := trajectorySegment.Length()

	type intersectingSegmentWrapper struct {
		point   vector.Vector2
		segment vector.Segment2
	}

	polyLen := len(poly)

	for i, _ := range poly {
		p1 := poly[i]
		p2 := poly[(i+1)%polyLen]

		// Il s'agit bien d'une intersection de lignes et pas de segments
		// Car le collider peut entre en collision avec la ligne formée par le segment même si son centre n'entre pas en collision avec le segment (le radius du collider est non nul)
		points = append(points, LineCircleIntersectionPoints(p1, p2, center, radius)...)
	}

	trajectorySlope, _ /*colliderAffineYIntersect*/, trajectoryIsVertical, _ /*colliderAffineVerticalX*/ := GetAffineEquationExpressedForY(trajectorySegment)
	if trajectoryIsVertical {
		trajectorySegment = trajectorySegment.SetPointB(trajectorySegment.GetPointB().Add(vector.MakeVector2(0.0001, 0)))
		trajectorySlope, _ /*YIntersect*/, trajectoryIsVertical, _ /*VerticalX*/ = GetAffineEquationExpressedForY(trajectorySegment)
		if trajectoryIsVertical {
			// no collision will be processed !
			// may never happen
			return []vector.Segment2{}
		}
	}

	orthoSlope := -1 / trajectorySlope //perpendicular(y=ax+b) : y = -1/a

	centerLinePoints := make([]vector.Vector2, 0)

	if len(points) == 0 {
		// Pas d'intersection de la surface de trajectoire du collider avec celle du cercle de collidee
		// Le cercle est peut-être trop petit

		if PointIsInTriangle(center, poly[0], poly[1], poly[2]) || PointIsInTriangle(center, poly[2], poly[3], poly[0]) {

			// on projette orthogonalement le centre du cercle sur la trajectoire du centre
			// on détermine l'orthogonale à la ligne de centre passant par le centre du cercle

			// on connaît un premier point sur l'orthogonale; on en cherche un deuxième
			// on détermine un deuxième point de la ligne
			// y1 = (m*x1) + (y0 - m*x0)

			orthoPrimePointX := center.GetX() + 10
			orthoPrimePointY := orthoSlope*orthoPrimePointX + (center.GetY() - orthoSlope*center.GetX())

			// on détermine le point d'intersection de l'orthogonale sur la ligne de trajectoire
			orthoCenterPoint, _ := LinesIntersectionPoint(
				trajectorySegment.GetPointA(), trajectorySegment.GetPointB(),
				center, vector.MakeVector2(orthoPrimePointX, orthoPrimePointY),
			)

			centerLinePoints = append(centerLinePoints, orthoCenterPoint)
		}
	} else {

		// Pour chaque point, on trouve son intersection avec la ligne centrale de trajectoire
		for _, p := range points {
			// on détermine l'orthogonale à la ligne de centre passant par le point en question

			// on connaît un premier point sur l'orthogonale; on en cherche un deuxième
			//orthoYIntersect := p.GetY() - (orthoSlope * p.GetX())

			// on détermine un deuxième point de la ligne
			// y1 = (m*x1) + (y0 - m*x0)

			orthoPrimePointX := p.GetX() + 10
			orthoPrimePointY := orthoSlope*orthoPrimePointX + (p.GetY() - orthoSlope*p.GetX())

			// on détermine le point d'intersection de l'orthogonale sur la ligne de trajectoire
			orthoCenterPoint, _ := LinesIntersectionPoint(
				trajectorySegment.GetPointA(), trajectorySegment.GetPointB(),
				p, vector.MakeVector2(orthoPrimePointX, orthoPrimePointY),
			)

			centerLinePoints = append(centerLinePoints, orthoCenterPoint)
		}
	}

	if len(centerLinePoints) >= 2 {
		// on identifie les distances min et max de projection des collisions sur la ligne de centre
		minDist := -1.0
		maxDist := -1.0
		for _, centerLinePoint := range centerLinePoints {
			dist := centerLinePoint.Sub(trajectorySegment.GetPointA()).Mag()
			if minDist < 0 || dist < minDist {
				minDist = dist
			}

			if maxDist < 0 || dist > maxDist {
				maxDist = dist
			}
		}

		if minDist > maxDist {
			minDist, maxDist = maxDist, minDist
		}

		minDist = minDist - colliderRadius
		maxDist = maxDist + colliderRadius

		if minDist < 0 {
			minDist = 0.001
		}

		if maxDist > trajectoryLength {
			maxDist = trajectoryLength
		}

		if minDist > maxDist {
			minDist, maxDist = maxDist, minDist
		}

		//log.Println("minDist", "maxDist", minDist, maxDist)

		collisionPositionSegments = append(collisionPositionSegments, trajectorySegment.SetLengthFromA(minDist))
		collisionPositionSegments = append(collisionPositionSegments, trajectorySegment.SetLengthFromA(maxDist))
	}

	return collisionPositionSegments
}
