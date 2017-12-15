package trigo

import (
	"math"

	"github.com/bytearena/core/common/utils/number"
	"github.com/bytearena/core/common/utils/vector"
)

func collideConstrainedCenterCircleWithPolygon(crossingPoly []vector.Vector2, centerSegment vector.Segment2, circleRadius float64) []vector.Segment2 {

	tangentsRadiuses := make([]vector.Segment2, 0) // point A: center position of object when colliding; point B: collision point

	/*
		TODO:
			collision du cercle du collider, dont le centre est sur la ligne centrale de trajectoire de l'objet
			avec le polygone obtenu par le clipping du rectangle orienté de trajectoire du collider avec le rectangle orienté de trajectoire du collidee.

			Les points de collision correspondent aux centres des cercles (aux positions des objets) tangents au début et à la fin du polygone.

			Pour identifier les faces du polygone sur lesquelles déterminer une tangente, il faut vérifier si la face du polygone testée est parallèle la ligne de centre de trajectoire de l'objet;
			si c'est le cas, il ne faut pas déterminer de tangente pour la face en question (utiliser pour ce test trigo.IntersectionWithLineSegment())
	*/

	// 1. On biaise les droites exactement verticales pour pouvoir toujours les décrire avec une équation affine
	colliderAffineSlope, _ /*colliderAffineYIntersect*/, colliderAffineIsVertical, _ /*colliderAffineVerticalX*/ := GetAffineEquationExpressedForY(centerSegment)
	if colliderAffineIsVertical {
		centerSegment = centerSegment.SetPointB(centerSegment.GetPointB().Add(vector.MakeVector2(0.0001, 0)))
		colliderAffineSlope, _ /*colliderAffineYIntersect*/, colliderAffineIsVertical, _ /*colliderAffineVerticalX*/ = GetAffineEquationExpressedForY(centerSegment)
		if colliderAffineIsVertical {
			// no collision will be processed !
			// may never happen
			return tangentsRadiuses
		}
		//panic("colliderAffineIsVertical !! what should we do ???")
	}

	// 2. On détermine les segments du polygone pour lesquels calculer une tangente

	type intersectingSegmentWrapper struct {
		point   vector.Vector2
		segment vector.Segment2
	}

	polyLen := len(crossingPoly)

	centerSegmentPointA, centerSegmentPointB := centerSegment.Get()
	touchingSegments := make([]intersectingSegmentWrapper, 0)

	for i, _ := range crossingPoly {
		p1 := crossingPoly[i]
		p2 := crossingPoly[(i+1)%polyLen]

		// Il s'agit bien d'une intersection de lignes et pas de segments
		// Car le collider peut entre en collision avec la ligne formée par le segment même si son centre n'entre pas en collision avec le segment (le radius du collider est non nul)
		if point, parallel := LinesIntersectionPoint(p1, p2, centerSegmentPointA, centerSegmentPointB); !parallel {
			touchingSegments = append(touchingSegments, intersectingSegmentWrapper{
				point:   point,
				segment: vector.MakeSegment2(p1, p2),
			})
		}
	}

	if len(touchingSegments) == 0 {
		// Pas d'intersection de la surface de trajectoire du collider avec celle du collidee
		// Ne devrait pas se produire, car ce cas est rendu impossible par le test 1, et par le fait que le polygone est une intersection des deux trajectoires considérées
		return tangentsRadiuses
	}

	// collision du cercle du collider, dont le centre est sur la ligne centrale de trajectoire de l'objet
	// avec le polygone obtenu par le clipping du rectangle orienté de trajectoire du collider avec le rectangle orienté de trajectoire du collidee.
	for _, touchingSegment := range touchingSegments {
		// On détermine l'équation affine de la droite passant entre les deux points du segment
		segmentAffineSlope, _ /*segmentAffineYIntersect*/, segmentAffineIsVertical, _ /*segmentAffineVerticalX*/ := GetAffineEquationExpressedForY(touchingSegment.segment)
		if segmentAffineIsVertical {
			//panic("segmentAffineIsVertical !! what should we do ???")
			touchingSegment.segment = touchingSegment.segment.SetPointB(touchingSegment.segment.GetPointB().Add(vector.MakeVector2(0.0001, 0)))
			segmentAffineSlope, _ /*colliderAffineYIntersect*/, segmentAffineIsVertical, _ /*colliderAffineVerticalX*/ = GetAffineEquationExpressedForY(touchingSegment.segment)
			if segmentAffineIsVertical {
				// no collision will be processed !
				// may never happen
				continue
			}
		}

		// le rayon pour la tangente recherchée (rt) est perpendiculaire au segment
		tangentRadiusSlope := -1 / segmentAffineSlope //perpendicular(y=ax+b) : y = -1/a
		var tangentRadiusSegment vector.Segment2

		// le centre du rayon (h, k) pour la tangente recherchée est le point d'intersection de rt et de la ligne de centre du collider
		if number.FloatEquals(tangentRadiusSlope, colliderAffineSlope) {
			// si le tangentRadiusSlope == colliderAffineSlope, la ligne de centre du collider est perpendiculaire au segment
			// le rayon au point de tangente est colinéaire à la ligne de centre du collider

			// On crée le segment depuis le début de la ligne de centre du collider jusqu'à sa collision avec la ligne (pas le segment) formée par le segment collidee
			tangentRadiusSegment = vector.MakeSegment2(centerSegmentPointA, touchingSegment.point).SetLengthFromB(circleRadius)
		} else {
			// il faut calculer le point d'intersection du segment perpendiculaire au collidee sur la ligne de centre du collider, et de longueur circleRadius
			// Utilisation du théorème de pythagore pour ce faire
			/*

					|\
				a	| \  c
					|  \
					|___\
					  b

					  c: ligne du collider
					  a: ligne du collidee
					  b: rayon de la tangente au cercle du collider

					On veut déterminer a
					On connaît b (circleRadius)
					On connaît la slope de l'angle ab (perpendiculaire)
					Il faut calculer l'angle ac
					Utiliser cet angle pour déterminer la slope (relation entre les longueurs a et c)
					Utiliser b, la slope de l'angle ac et le fait que le triangle soit rectangle ab pour calculer a

					Comme le triangle est rectangle, slopeac = a/b
					Donc: a = slopeac * b

					On utilise a pour déterminer b, et on recule depuis le point de collision de la ligne du collider de la longueur déterminée pour trouver le centre du cercle tangent
			*/

			absoluteAngleCRad := math.Atan2(
				centerSegmentPointA.GetY()-centerSegmentPointB.GetY(),
				centerSegmentPointA.GetX()-centerSegmentPointB.GetX(),
			)
			absoluteAngleARad := math.Atan2(
				touchingSegment.segment.GetPointA().GetY()-touchingSegment.segment.GetPointB().GetY(),
				touchingSegment.segment.GetPointA().GetX()-touchingSegment.segment.GetPointB().GetX(),
			)

			angleACRad := absoluteAngleCRad - absoluteAngleARad
			slopeAC := math.Tan(angleACRad)

			a := slopeAC * circleRadius // la distance entre le point d'intersection de la ligne du collider et la tangente du cercle de l'agent

			// on utilise a et b pour déterminer c
			// a2 + b2 = c2
			// c = sqrt(a2+b2)

			c := math.Sqrt(math.Pow(a, 2) + math.Pow(circleRadius, 2))

			// On recule depuis le point de collision de la ligne du collider de la longueur déterminée pour trouver le centre du cercle tangent
			colliderCollisionCourseSegment := vector.MakeSegment2(centerSegmentPointA, touchingSegment.point).SetLengthFromB(c)
			tangentCircleCenter := colliderCollisionCourseSegment.GetPointA()

			// on calcule l'intersection sur collidee du rayon de tangente passant en tangentCircleCenter
			// on connait le slope du rayon de tangente tangentRadiusSlope
			// on connait un point par lequel passe ce rayon tangentCircleCenter
			// on cherche l'équation affine de la ligne en question
			// y0 - y1 = m(x0 - x1)
			// y0 - y1 = m*x0 - m*x1
			// -y1 = m*x0 - m*x1 - y0
			// y1 = -1 * (m*x0 - m*x1 - y0)
			// y1 = (m*x1) + (y0 - m*x0)
			// y=ax+b
			// a=m
			// b=y0-m*x0

			tangentRadiusAffineYIntersect := tangentCircleCenter.GetY() - (tangentRadiusSlope * tangentCircleCenter.GetX())
			// on détermine un deuxième point de la ligne
			tangentRadiusPrimePointX := tangentCircleCenter.GetX() + 10
			tangentRadiusPrimePointY := tangentRadiusSlope*tangentRadiusPrimePointX + tangentRadiusAffineYIntersect

			// on détermine le point d'intersection du rayon de tangente
			tangentPoint, _ := LinesIntersectionPoint(
				touchingSegment.segment.GetPointA(), touchingSegment.segment.GetPointB(),
				tangentCircleCenter, vector.MakeVector2(tangentRadiusPrimePointX, tangentRadiusPrimePointY),
			)

			tangentRadiusSegment = vector.MakeSegment2(tangentCircleCenter, tangentPoint)
		}

		tangentsRadiuses = append(tangentsRadiuses, tangentRadiusSegment)

	}

	// Les points de collision correspondent aux centres des cercles (aux positions des objets) tangents au début et à la fin du polygone.

	return tangentsRadiuses
}
