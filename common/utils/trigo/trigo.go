package trigo

import (
	"errors"
	"math"

	"github.com/bytearena/core/common/utils/number"
	"github.com/bytearena/core/common/utils/vector"
)

func SegmentSegmentIntersection(p, q vector.Segment2) (intersection vector.Vector2, intersects bool, colinear bool, parallel bool) {
	p1, p2 := p.Get()
	q1, q2 := q.Get()
	return IntersectionWithLineSegment(p1, p2, q1, q2)
}

func IntersectionWithLineSegment(p vector.Vector2, p2 vector.Vector2, q vector.Vector2, q2 vector.Vector2) (intersection vector.Vector2, intersects bool, colinear bool, parallel bool) {

	r := p2.Sub(p)
	s := q2.Sub(q)
	rxs := r.Cross(s)
	qpxr := q.Sub(p).Cross(r)

	// If r x s = 0 and (q - p) x r = 0, then the two lines are collinear.
	if number.IsZero(rxs) && number.IsZero(qpxr) {
		// 1. If either  0 <= (q - p) * r <= r * r or 0 <= (p - q) * s <= * s
		// then the two lines are overlapping,
		qSubPTimesR := q.Sub(p).Dot(r)
		pSubQTimesS := p.Sub(q).Dot(s)
		rSquared := r.Dot(r)
		sSquared := s.Dot(s)

		if (qSubPTimesR >= 0 && qSubPTimesR <= rSquared) || (pSubQTimesS >= 0 && pSubQTimesS <= sSquared) {
			return vector.MakeNullVector2(), true, true, true
		}

		// 2. If neither 0 <= (q - p) * r = r * r nor 0 <= (p - q) * s <= s * s
		// then the two lines are collinear but disjoint.
		// No need to implement this expression, as it follows from the expression above.
		return vector.MakeNullVector2(), false, true, true
	}

	// 3. If r x s = 0 and (q - p) x r != 0, then the two lines are parallel and non-intersecting.
	if number.IsZero(rxs) && !number.IsZero(qpxr) {
		return vector.MakeNullVector2(), false, false, true
	}

	t := q.Sub(p).Cross(s) / rxs
	u := q.Sub(p).Cross(r) / rxs

	// 4. If r x s != 0 and 0 <= t <= 1 and 0 <= u <= 1
	// the two line segments meet at the point p + t r = q + u s.
	if !number.IsZero(rxs) && (0 <= t && t <= 1) && (0 <= u && u <= 1) {
		// We can calculate the intersection point using either t or u.
		return p.Add(r.MultScalar(t)), true, false, false
	}

	// 5. Otherwise, the two line segments are not parallel but do not intersect.
	return vector.MakeNullVector2(), false, false, true
}

func IntersectionWithLineSegmentCheckOnly(p1 vector.Vector2, p2 vector.Vector2, p3 vector.Vector2, p4 vector.Vector2) (intersect bool) {
	a := p2.Sub(p1)
	b := p3.Sub(p4)
	c := p1.Sub(p3)

	ax, ay := a.Get()
	bx, by := b.Get()
	cx, cy := c.Get()

	alphaNumerator := by*cx - bx*cy
	alphaDenominator := ay*bx - ax*by
	betaNumerator := ax*cy - ay*cx
	betaDenominator := alphaDenominator

	doIntersect := true

	if alphaDenominator == 0 || betaDenominator == 0 {
		doIntersect = false
	} else {
		if alphaDenominator > 0 {
			if alphaNumerator < 0 || alphaNumerator > alphaDenominator {
				doIntersect = false
			}
		} else if alphaNumerator > 0 || alphaNumerator < alphaDenominator {
			doIntersect = false
		}

		if doIntersect && betaDenominator > 0 {
			if betaNumerator < 0 || betaNumerator > betaDenominator {
				doIntersect = false
			}
		} else if betaNumerator > 0 || betaNumerator < betaDenominator {
			doIntersect = false
		}
	}

	return doIntersect
}

func LinesIntersectionPoint(p0 vector.Vector2, p1 vector.Vector2, p2 vector.Vector2, p3 vector.Vector2) (point vector.Vector2, parallel bool) {

	p0x, p0y := p0.Get()
	p1x, p1y := p1.Get()
	p2x, p2y := p2.Get()
	p3x, p3y := p3.Get()

	s1_x := p1x - p0x
	s1_y := p1y - p0y
	s2_x := p3x - p2x
	s2_y := p3y - p2y

	s := (-s1_y*(p0x-p2x) + s1_x*(p0y-p2y)) / (-s2_x*s1_y + s1_x*s2_y)
	t := (s2_x*(p0y-p2y) - s2_y*(p0x-p2x)) / (-s2_x*s1_y + s1_x*s2_y)

	if s >= 0 && s <= 1 && t >= 0 && t <= 1 {
		// Collision detected
		return vector.MakeVector2(p0x+(t*s1_x), p0y+(t*s1_y)), false
	}

	// No collision
	return vector.MakeNullVector2(), true
}

// http://devmag.org.za/2009/04/17/basic-collision-detection-in-2d-part-2/
func LineCircleIntersectionPoints(LineP1 vector.Vector2, LineP2 vector.Vector2, CircleCentre vector.Vector2, Radius float64) []vector.Vector2 {

	LocalP1 := LineP1.Sub(CircleCentre)
	LocalP2 := LineP2.Sub(CircleCentre)
	// Precalculate this value. We use it often
	P2MinusP1 := LocalP2.Sub(LocalP1)

	p2minusp1x, p2minusp1y := P2MinusP1.Get()
	localp1x, localp1y := LocalP1.Get()

	a := P2MinusP1.MagSq()
	b := 2 * ((p2minusp1x * localp1x) + (p2minusp1y * localp1y))
	c := LocalP1.MagSq() - (Radius * Radius)

	delta := b*b - (4 * a * c)
	if delta < 0 {
		// No intersection
		return make([]vector.Vector2, 0)
	}

	if delta == 0 {
		u := -b / (2.0 * a)

		// Use LineP1 instead of LocalP1 because we want our answer in global space, not the circle's local space
		res := make([]vector.Vector2, 1)
		res[0] = LineP1.Add(P2MinusP1.MultScalar(u))
		return res
	}

	// (delta > 0) // Two intersections
	SquareRootDelta := math.Sqrt(delta)

	u1 := (-b + SquareRootDelta) / (2 * a)
	u2 := (-b - SquareRootDelta) / (2 * a)

	res := make([]vector.Vector2, 2)
	res[0] = LineP1.Add(P2MinusP1.MultScalar(u1))
	res[1] = LineP1.Add(P2MinusP1.MultScalar(u2))

	return res
}

// Taken from https://stackoverflow.com/a/12221389
// Initially from Tim Voght, http://paulbourke.net/geometry/circlesphere/tvoght.c
func CircleCircleIntersectionPoints(center0 vector.Vector2, radius0 float64, center1 vector.Vector2, radius1 float64) (intersections []vector.Vector2, firstContainsSecond bool, secondContainsFirst bool) {
	// var a, dx, dy, d, h, rx, ry;
	// var x2, y2;

	x0, y0 := center0.Get()
	x1, y1 := center1.Get()

	/* dx and dy are the vertical and horizontal distances between
	 * the circle centers.
	 */
	dx := x1 - x0
	dy := y1 - y0

	/* Determine the straight-line distance between the centers. */
	d := math.Sqrt((dy * dy) + (dx * dx))

	/* Check for solvability. */
	if d > (radius0 + radius1) {
		/* no solution. circles do not intersect. */
		return []vector.Vector2{}, false, false
	}

	if d < math.Abs(radius0-radius1) {
		/* no solution. one circle is contained in the other */
		if radius0 > radius1 {
			// first contains second
			return []vector.Vector2{}, true, false
		}

		// second contains first
		return []vector.Vector2{}, false, true
	}

	/* 'point 2' is the point where the line through the circle
	 * intersection points crosses the line between the circle
	 * centers.
	 */

	/* Determine the distance from point 0 to point 2. */
	a := (math.Pow(radius0, 2) - math.Pow(radius1, 2) + math.Pow(d, 2)) / (2.0 * d)

	/* Determine the coordinates of point 2. */
	x2 := x0 + (dx * a / d)
	y2 := y0 + (dy * a / d)

	/* Determine the distance from point 2 to either of the
	 * intersection points.
	 */
	h := math.Sqrt(math.Pow(radius0, 2) - math.Pow(a, 2))

	/* Now determine the offsets of the intersection points from
	 * point 2.
	 */
	rx := -dy * (h / d)
	ry := dx * (h / d)

	/* Determine the absolute intersection points. */
	xi := x2 + rx
	xiPrime := x2 - rx
	yi := y2 + ry
	yiPrime := y2 - ry

	return []vector.Vector2{vector.MakeVector2(xi, yi), vector.MakeVector2(xiPrime, yiPrime)}, false, false
}

func PointOnLineSegment(p vector.Vector2, a vector.Vector2, b vector.Vector2) bool {
	t := 0.0001

	px, py := p.Get()
	ax, ay := a.Get()
	bx, by := b.Get()

	// ensure points are collinear
	zero := (bx-ax)*(py-ay) - (px-ax)*(by-ay)
	if zero > t || zero < -t {
		return false
	}

	// check if x-coordinates are not equal
	if ax-bx > t || bx-ax > t {
		// ensure x is between a.x & b.x (use tolerance)
		if ax > bx {
			return px+t > bx && px-t < ax
		} else {
			return px+t > ax && px-t < bx
		}
	}

	// ensure y is between a.y & b.y (use tolerance)
	if ay > by {
		return py+t > by && py-t < ay
	}

	return py+t > ay && py-t < by
}

func FullCircleAngleToSignedHalfCircleAngle(rad float64) float64 {
	if rad > math.Pi { // 180° en radians
		rad -= math.Pi * 2 // 360° en radian
	} else if rad < -math.Pi {
		rad += math.Pi * 2 // 360° en radian
	}

	return rad
}

func PointIsInTriangle(point, p0, p1, p2 vector.Vector2) bool {
	p0x, p0y := p0.Get()
	p1x, p1y := p1.Get()
	p2x, p2y := p2.Get()
	px, py := point.Get()

	Area := 0.5 * (-p1y*p2x + p0y*(-p1x+p2x) + p0x*(p1y-p2y) + p1x*p2y)

	s := 1 / (2 * Area) * (p0y*p2x - p0x*p2y + (p2y-p0y)*px + (p0x-p2x)*py)
	t := 1 / (2 * Area) * (p0x*p1y - p0y*p1x + (p0y-p1y)*px + (p1x-p0x)*py)

	return (s > 0 && t > 0 && 1-s-t > 0)
}

func PointIsInCircle(point vector.Vector2, center vector.Vector2, radius float64) bool {
	squareDist := math.Pow(center.GetX()-point.GetX(), 2) + math.Pow(center.GetY()-point.GetY(), 2)
	return squareDist <= math.Pow(radius, 2)
}

func ComputeCenterOfMass(points []vector.Vector2) (vector.Vector2, error) {
	if len(points) == 0 {
		return vector.MakeNullVector2(), errors.New("Cannot compute center of mass on empty list")
	}

	sumx, sumy := 0.0, 0.0
	for _, p := range points {
		sumx += p.GetX()
		sumy += p.GetY()
	}

	flen := float64(len(points))
	return vector.MakeVector2(sumx/flen, sumy/flen), nil
}

/**
 * Helper function to determine whether there is an intersection between the two polygons described
 * by the lists of vertices. Uses the Separating Axis Theorem
 *
 * @param a an array of connected points [{x:, y:}, {x:, y:},...] that form a closed polygon
 * @param b an array of connected points [{x:, y:}, {x:, y:},...] that form a closed polygon
 * @return true if there is any intersection between the 2 polygons, false otherwise
 */
// taken from https://stackoverflow.com/a/12414951

func DoClosedConvexPolygonsIntersect(a []vector.Vector2, b []vector.Vector2) bool {
	polygons := []([]vector.Vector2){a, b}

	var minA *float64 /* = nil*/
	var maxA *float64 /* = nil*/
	var minB *float64 /* = nil*/
	var maxB *float64 /* = nil*/

	for _, polygon := range polygons {

		// for each polygon, look at each edge of the polygon, and determine if it separates
		// the two shapes

		polylen := len(polygon)

		for i1 := 0; i1 < len(polygon); i1++ {
			// grab 2 vertices to create an edge

			i2 := (i1 + 1) % polylen

			p1 := polygon[i1]
			p2 := polygon[i2]

			// find the line perpendicular to this edge
			normalX := p2.GetY() - p1.GetY()
			normalY := p1.GetX() - p2.GetX()

			minA = nil
			maxA = nil

			// for each vertex in the first shape, project it onto the line perpendicular to the edge
			// and keep track of the min and max of these values

			for j := 0; j < len(a); j++ {
				projected := normalX*a[j].GetX() + normalY*a[j].GetY()
				if minA == nil || projected < *minA {
					minA = &projected
				}

				if maxA == nil || projected > *maxA {
					maxA = &projected
				}
			}

			// for each vertex in the second shape, project it onto the line perpendicular to the edge
			// and keep track of the min and max of these values

			minB = nil
			maxB = nil

			for j := 0; j < len(b); j++ {
				projected := normalX*b[j].GetX() + normalY*b[j].GetY()
				if minB == nil || projected < *minB {
					minB = &projected
				}

				if maxB == nil || projected > *maxB {
					maxB = &projected
				}
			}

			// if there is no overlap between the projects, the edge we are looking at separates the two
			// polygons, and we know there is no overlap
			if *maxA < *minB || *maxB < *minA {
				return false
			}
		}
	}

	return true
}

func GetAffineEquationExpressedForY(segment vector.Segment2) (a float64, b float64, vertical bool, xvertical float64) { // y = ax + b
	pA, pB := segment.Get()
	pAX, pAY := pA.Get()
	pBX, pBY := pB.Get()

	if number.IsZero(pBX - pAX) {
		// line is vertical
		return 0, 0, true, pAX
	}

	a = (pBY - pAY) / (pBX - pAX)
	b = pAY - (a * pAX)
	return a, b, false, 0
}

func LocalAngleToAbsoluteAngleVec(abscurrentagentangle float64, vec vector.Vector2, maxangleconstraint *float64) vector.Vector2 {

	// On passe de 0° / 360° à -180° / +180°
	relvecangle := FullCircleAngleToSignedHalfCircleAngle(vec.Angle())

	// On contraint la vélocité angulaire à un maximum
	if maxangleconstraint != nil {
		maxangleconstraintval := *maxangleconstraint
		if math.Abs(relvecangle) > maxangleconstraintval {
			if relvecangle > 0 {
				relvecangle = maxangleconstraintval
			} else {
				relvecangle = -1 * maxangleconstraintval
			}
		}
	}

	return vec.SetAngle(abscurrentagentangle + relvecangle)
}

func GetBoundingBoxForPoints(points []vector.Vector2) (lowerBound vector.Vector2, upperBound vector.Vector2) {

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

	width := maxX - minX
	if width <= 0 {
		width = 0.00001
	}

	height := maxY - minY
	if height <= 0 {
		height = 0.00001
	}

	return vector.MakeVector2(minX, minY), vector.MakeVector2(maxX, maxY)
}
