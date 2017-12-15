package trigo

import (
	"github.com/bytearena/core/common/utils/vector"
)

func collideCirclesCircles(colliderCenterA vector.Vector2, colliderRadiusA float64, collideeCenterA vector.Vector2, collideeRadiusA float64, colliderCenterB vector.Vector2, colliderRadiusB float64, collideeCenterB vector.Vector2, collideeRadiusB float64) []vector.Vector2 {

	points := make([]vector.Vector2, 0)

	if colliderCenterA.Equals(collideeCenterA) || colliderCenterA.Equals(collideeCenterB) {
		return []vector.Vector2{colliderCenterA}
	}

	if colliderCenterB.Equals(collideeCenterA) || colliderCenterB.Equals(collideeCenterB) {
		return []vector.Vector2{colliderCenterB}
	}

	// ColliderCircleA/CollideeCircleA, ColliderCircleB/CollideeCircleB
	// ColliderCircleA/CollideeCircleB, ColliderCircleB/CollideeCircleA

	// ColliderCircleA/CollideeCircleA
	intersections, firstContainsSecond, secondContainsFirst := CircleCircleIntersectionPoints(colliderCenterA, colliderRadiusA, collideeCenterA, collideeRadiusA)
	if len(intersections) > 0 {
		points = append(points, intersections...)
	} else if firstContainsSecond {
		points = append(points, collideeCenterA)
	} else if secondContainsFirst {
		points = append(points, colliderCenterA)
	}

	// ColliderCircleB/CollideeCircleB
	intersections, firstContainsSecond, secondContainsFirst = CircleCircleIntersectionPoints(colliderCenterB, colliderRadiusB, collideeCenterB, collideeRadiusB)
	if len(intersections) > 0 {
		points = append(points, intersections...)
	} else if firstContainsSecond {
		points = append(points, collideeCenterB)
	} else if secondContainsFirst {
		points = append(points, colliderCenterB)
	}

	// ColliderCircleA/CollideeCircleB
	intersections, firstContainsSecond, secondContainsFirst = CircleCircleIntersectionPoints(colliderCenterA, colliderRadiusA, collideeCenterB, collideeRadiusB)
	if len(intersections) > 0 {
		points = append(points, intersections...)
	} else if firstContainsSecond {
		points = append(points, collideeCenterB)
	} else if secondContainsFirst {
		points = append(points, colliderCenterA)
	}

	// ColliderCircleB/CollideeCircleA
	intersections, firstContainsSecond, secondContainsFirst = CircleCircleIntersectionPoints(colliderCenterB, colliderRadiusB, collideeCenterA, collideeRadiusA)
	if len(intersections) > 0 {
		points = append(points, intersections...)
	} else if firstContainsSecond {
		points = append(points, collideeCenterA)
	} else if secondContainsFirst {
		points = append(points, colliderCenterB)
	}

	return points
}
