package visibility2d

import (
	"math"
	"sort"

	"github.com/bytearena/core/common/types/datastructures"
)

func OnlyVisible(position [2]float64, perceptionitems []ObstacleSegment) []ObstacleSegment {

	rescaleFactor := 10000.0

	scaledUpSegments := make([]ObstacleSegment, 0)
	for _, item := range perceptionitems {
		scaledUpSegments = append(scaledUpSegments, ObstacleSegment{
			Points: [2][2]float64{
				[2]float64{item.Points[0][0] * rescaleFactor, item.Points[0][1] * rescaleFactor},
				[2]float64{item.Points[1][0] * rescaleFactor, item.Points[1][1] * rescaleFactor},
			},
			UserData: item.UserData,
		})
	}

	scaledUpSegments = breakIntersections(scaledUpSegments)

	//brokenSegments := breakIntersections(scaledUpSegments)

	visibleSegments := make([]ObstacleSegment, 0)

	visibility := makeVisibilityProcessor()

	for _, item := range scaledUpSegments {
		visibility.AddSegment(
			item.Points[0][0], item.Points[0][1],
			item.Points[1][0], item.Points[1][1],
			item.UserData,
		)
	}

	visibility.SetLightLocation(position[0], position[1])
	visibility.Sweep()

	for _, visibleSegment := range visibility.output {
		visibleSegments = append(visibleSegments, ObstacleSegment{
			Points: [2][2]float64{
				[2]float64{visibleSegment.p1[0] / rescaleFactor, visibleSegment.p1[1] / rescaleFactor},
				[2]float64{visibleSegment.p2[0] / rescaleFactor, visibleSegment.p2[1] / rescaleFactor},
			},
			UserData: visibleSegment.completeSegment.userData,
		})
	}

	return visibleSegments
}

type point [2]float64

type endPoint struct {
	point
	begin   bool
	segment *visibilityComputationSegment
	angle   float64
}

type visibilityComputationSegment struct {
	p1       *endPoint
	p2       *endPoint
	d        float64
	userData interface{}
}

type visibleSegment struct {
	p1              point
	p2              point
	completeSegment *visibilityComputationSegment
}

type visibilityProcessor struct {
	segments  []*visibilityComputationSegment
	endPoints []*endPoint
	center    point
	open      datastructures.DLL
	output    []visibleSegment
	//intersectionsDetected [][]point
}

func makeVisibilityProcessor() visibilityProcessor {
	return visibilityProcessor{
		segments:  make([]*visibilityComputationSegment, 0),
		endPoints: make([]*endPoint, 0),
		open:      datastructures.DLL{},
		center:    point{0, 0},
		output:    make([]visibleSegment, 0),
		//intersectionsDetected: make([][]point, 0),
	}
}

func (visi *visibilityProcessor) AddSegment(x1, y1, x2, y2 float64, userdata interface{}) {

	seg := &visibilityComputationSegment{}

	p1 := &endPoint{
		point:   point{x1, y1},
		segment: seg,
	}

	p2 := &endPoint{
		point:   point{x2, y2},
		segment: seg,
	}

	seg.p1 = p1
	seg.p2 = p2
	seg.d = 0.0
	seg.userData = userdata

	visi.segments = append(visi.segments, seg)
	visi.endPoints = append(visi.endPoints, p1, p2)

}

func (visi *visibilityProcessor) SetLightLocation(x, y float64) {
	visi.center[0] = x
	visi.center[1] = y

	for _, seg := range visi.segments {

		dx := 0.5*(seg.p1.point[0]+seg.p2.point[0]) - x
		dy := 0.5*(seg.p1.point[1]+seg.p2.point[1]) - y
		// NOTE: we only use this for comparison so we can use
		// distance squared instead of distance. However in
		// practice the sqrt is plenty fast and this doesn't
		// really help in this situation.
		seg.d = dx*dx + dy*dy

		// NOTE: future optimization: we could record the quadrant
		// and the y/x or x/y ratio, and sort by (quadrant,
		// ratio), instead of calling atan2. See
		// <https://github.com/mikolalysenko/compare-slope> for a
		// library that does this. Alternatively, calculate the
		// angles and use bucket sort to get an O(N) sort.
		seg.p1.angle = math.Atan2(seg.p1.point[1]-y, seg.p1.point[0]-x)
		seg.p2.angle = math.Atan2(seg.p2.point[1]-y, seg.p2.point[0]-x)

		dAngle := seg.p2.angle - seg.p1.angle
		if dAngle <= -math.Pi {
			dAngle += 2 * math.Pi
		}
		if dAngle > math.Pi {
			dAngle -= 2 * math.Pi
		}
		seg.p1.begin = (dAngle > 0.0)
		seg.p2.begin = !seg.p1.begin
	}
}

type byAngle []*endPoint

func (coll byAngle) Len() int      { return len(coll) }
func (coll byAngle) Swap(i, j int) { coll[i], coll[j] = coll[j], coll[i] }
func (coll byAngle) Less(i, j int) bool {

	a := coll[i]
	b := coll[j]

	// Traverse in angle order
	if a.angle > b.angle {
		return false
	}

	if a.angle < b.angle {
		return true
	}

	// But for ties (common), we want Begin nodes before End nodes
	if !a.begin && b.begin {
		return false
	}

	if a.begin && !b.begin {
		return true
	}

	return false
}

func leftOf(s *visibilityComputationSegment, p point) bool {
	// This is based on a 3d cross product, but we don't need to
	// use z coordinate inputs (they're 0), and we only need the
	// sign. If you're annoyed that cross product is only defined
	// in 3d, see "outer product" in Geometric Algebra.
	// <http://en.wikipedia.org/wiki/Geometric_algebra>
	cross := (s.p2.point[0]-s.p1.point[0])*(p[1]-s.p1.point[1]) - (s.p2.point[1]-s.p1.point[1])*(p[0]-s.p1.point[0])
	return cross < 0
	// Also note that this is the naive version of the test and
	// isn't numerically robust. See
	// <https://github.com/mikolalysenko/robust-arithmetic> for a
	// demo of how this fails when a point is very close to the
	// line.
}

func interpolate(p, q point, f float64) point {
	return point{p[0]*(1-f) + q[0]*f, p[1]*(1-f) + q[1]*f}
}

// Helper: do we know that segment a is in front of b?
// Implementation not anti-symmetric (that is to say,
// _segment_in_front_of(a, b) != (!_segment_in_front_of(b, a)).
// Also note that it only has to work in a restricted set of cases
// in the visibility algorithm; I don't think it handles all
// cases. See http://www.redblobgames.com/articles/visibility/segment-sorting.html
func (visi *visibilityProcessor) _segment_in_front_of(a, b *visibilityComputationSegment, relativeTo point) bool {
	// NOTE: we slightly shorten the segments so that
	// intersections of the endpoints (common) don't count as
	// intersections in this algorithm
	var A1 = leftOf(a, interpolate(b.p1.point, b.p2.point, 0.00001))
	var A2 = leftOf(a, interpolate(b.p2.point, b.p1.point, 0.00001))
	var A3 = leftOf(a, relativeTo)
	var B1 = leftOf(b, interpolate(a.p1.point, a.p2.point, 0.00001))
	var B2 = leftOf(b, interpolate(a.p2.point, a.p1.point, 0.00001))
	var B3 = leftOf(b, relativeTo)

	// NOTE: this algorithm is probably worthy of a short article
	// but for now, draw it on paper to see how it works. Consider
	// the line A1-A2. If both B1 and B2 are on one side and
	// relativeTo is on the other side, then A is in between the
	// viewer and B. We can do the same with B1-B2: if A1 and A2
	// are on one side, and relativeTo is on the other side, then
	// B is in between the viewer and A.
	if B1 == B2 && B2 != B3 {
		return true
	}
	if A1 == A2 && A2 == A3 {
		return true
	}
	if A1 == A2 && A2 != A3 {
		return false
	}
	if B1 == B2 && B2 == B3 {
		return false
	}

	// If A1 != A2 and B1 != B2 then we have an intersection.
	// Expose it for the GUI to show a message. A more robust
	// implementation would split segments at intersections so
	// that part of the segment is in front and part is behind.
	// visi.intersectionsDetected = append(
	// 	visi.intersectionsDetected,
	// 	[]point{a.p1.point, a.p2.point, b.p1.point, b.p2.point},
	// )

	return false

	// NOTE: previous implementation was a.d < b.d. That's simpler
	// but trouble when the segments are of dissimilar sizes. If
	// you're on a grid and the segments are similarly sized, then
	// using distance will be a simpler and faster implementation.
}

// Run the algorithm, sweeping over all or part of the circle to find
// the visible area, represented as a set of triangles
func (visi *visibilityProcessor) Sweep() {
	maxAngle := 999.0

	visi.output = make([]visibleSegment, 0) // output set of triangles
	//visi.intersectionsDetected = make([][]point, 0)

	sort.Sort(byAngle(visi.endPoints))

	visi.open.Clear()
	var beginAngle = 0.0

	// At the beginning of the sweep we want to know which
	// segments are active. The simplest way to do this is to make
	// a pass collecting the segments, and make another pass to
	// both collect and process them. However it would be more
	// efficient to go through all the segments, figure out which
	// ones intersect the initial sweep line, and then sort them.
	for pass := 0; pass < 2; pass++ {
		for _, p := range visi.endPoints {
			if pass == 1 && p.angle > maxAngle {
				// Early exit for the visualization to show the sweep process
				break
			}

			var current_old *visibilityComputationSegment = nil
			if !visi.open.Empty() {
				current_old = visi.open.Head.Val.(*visibilityComputationSegment)
			}

			if p.begin {

				// Insert into the right place in the list
				node := visi.open.Head

				for node != nil {
					valAsSegment := node.Val.(*visibilityComputationSegment)

					if !visi._segment_in_front_of(p.segment, valAsSegment, visi.center) {
						break
					}

					node = node.Next
				}

				if node == nil {
					visi.open.Append(p.segment)
				} else {
					visi.open.InsertBefore(node, p.segment)
				}
			} else {
				visi.open.RemoveVal(p.segment)
			}

			var current_new *visibilityComputationSegment = nil
			if !visi.open.Empty() {
				current_new = visi.open.Head.Val.(*visibilityComputationSegment)
			}

			//log.Println(pass, current_old, current_new)

			if current_old != current_new {
				if pass == 1 {
					visi.addTriangle(beginAngle, p.angle, current_old)
				}
				beginAngle = p.angle
			}
		}
	}
}

func lineIntersection(p1, p2, p3, p4 point) point {
	// From http://paulbourke.net/geometry/lineline2d/
	var s = ((p4[0]-p3[0])*(p1[1]-p3[1]) - (p4[1]-p3[1])*(p1[0]-p3[0])) / ((p4[1]-p3[1])*(p2[0]-p1[0]) - (p4[0]-p3[0])*(p2[1]-p1[1]))
	return point{p1[0] + s*(p2[0]-p1[0]), p1[1] + s*(p2[1]-p1[1])}
}

func (visi *visibilityProcessor) addTriangle(angle1, angle2 float64, segment *visibilityComputationSegment) {

	if segment == nil {
		return
	}

	var p1 point = visi.center
	var p2 point = point{visi.center[0] + math.Cos(angle1), visi.center[1] + math.Sin(angle1)}
	var p3 point = point{0.0, 0.0}
	var p4 point = point{0.0, 0.0}

	// Stop the triangle at the intersecting segment
	p3[0] = segment.p1.point[0]
	p3[1] = segment.p1.point[1]
	p4[0] = segment.p2.point[0]
	p4[1] = segment.p2.point[1]

	var pBegin = lineIntersection(p3, p4, p1, p2)

	p2[0] = visi.center[0] + math.Cos(angle2)
	p2[1] = visi.center[1] + math.Sin(angle2)
	var pEnd = lineIntersection(p3, p4, p1, p2)

	visi.output = append(visi.output, visibleSegment{
		p1:              pBegin,
		p2:              pEnd,
		completeSegment: segment,
	})
}

////////////////////////////////////////////////////////////////////////////
// Break intersections /////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////

const epsilon = 0.000001

type ObstacleSegment struct {
	Points   [2][2]float64
	UserData interface{}
}

func distance(a, b [2]float64) float64 {
	dx := a[0] - b[0]
	dy := a[1] - b[1]
	return dx*dx + dy*dy
}

func equal(a, b [2]float64) bool {
	return math.Abs(a[0]-b[0]) < epsilon && math.Abs(a[1]-b[1]) < epsilon
}

func intersectLines(a1, a2, b1, b2 [2]float64) (p [2]float64, intersects bool) {
	var dbx = b2[0] - b1[0]
	var dby = b2[1] - b1[1]
	var dax = a2[0] - a1[0]
	var day = a2[1] - a1[1]

	var uB = dby*dax - dbx*day
	if uB != 0 {
		var ua = (dbx*(a1[1]-b1[1]) - dby*(a1[0]-b1[0])) / uB
		return [2]float64{a1[0] - ua*-dax, a1[1] - ua*-day}, true
	}

	return [2]float64{}, false
}

func isOnSegment(xi, yi, xj, yj, xk, yk float64) bool {
	return (xi <= xk || xj <= xk) && (xk <= xi || xk <= xj) &&
		(yi <= yk || yj <= yk) && (yk <= yi || yk <= yj)
}

func computeDirection(xi, yi, xj, yj, xk, yk float64) int {
	var a = (xk - xi) * (yj - yi)
	var b = (xj - xi) * (yk - yi)
	if a < b {
		return -1
	}

	if a > b {
		return 1
	}

	return 0
}

func doLineSegmentsIntersect(x1, y1, x2, y2, x3, y3, x4, y4 float64) bool {
	var d1 = computeDirection(x3, y3, x4, y4, x1, y1)
	var d2 = computeDirection(x3, y3, x4, y4, x2, y2)
	var d3 = computeDirection(x1, y1, x2, y2, x3, y3)
	var d4 = computeDirection(x1, y1, x2, y2, x4, y4)
	return (((d1 > 0 && d2 < 0) || (d1 < 0 && d2 > 0)) &&
		((d3 > 0 && d4 < 0) || (d3 < 0 && d4 > 0))) ||
		(d1 == 0 && isOnSegment(x3, y3, x4, y4, x1, y1)) ||
		(d2 == 0 && isOnSegment(x3, y3, x4, y4, x2, y2)) ||
		(d3 == 0 && isOnSegment(x1, y1, x2, y2, x3, y3)) ||
		(d4 == 0 && isOnSegment(x1, y1, x2, y2, x4, y4))
}

func breakIntersections(segments []ObstacleSegment) []ObstacleSegment {
	var output = make([]ObstacleSegment, 0)

	for i := 0; i < len(segments); i++ {
		intersections := make([][2]float64, 0)

		for j := 0; j < len(segments); j++ {

			if i == j {
				continue
			}

			if doLineSegmentsIntersect(segments[i].Points[0][0], segments[i].Points[0][1], segments[i].Points[1][0], segments[i].Points[1][1], segments[j].Points[0][0], segments[j].Points[0][1], segments[j].Points[1][0], segments[j].Points[1][1]) {

				if intersectPoint, intersects := intersectLines(segments[i].Points[0], segments[i].Points[1], segments[j].Points[0], segments[j].Points[1]); intersects {
					if equal(intersectPoint, segments[i].Points[0]) || equal(intersectPoint, segments[i].Points[1]) {
						continue
					}

					intersections = append(intersections, intersectPoint)
				}
			}
		}

		start := [2]float64{segments[i].Points[0][0], segments[i].Points[0][1]}

		for len(intersections) > 0 {

			var endIndex = 0
			var endDis = distance(start, intersections[0])

			for j := 1; j < len(intersections); j++ {
				var dis = distance(start, intersections[j])
				if dis < endDis {
					endDis = dis
					endIndex = j
				}
			}

			output = append(output, ObstacleSegment{
				Points: [2][2]float64{
					[2]float64{start[0], start[1]},
					[2]float64{intersections[endIndex][0], intersections[endIndex][1]},
				},
				UserData: segments[i].UserData,
			})
			start[0] = intersections[endIndex][0]
			start[1] = intersections[endIndex][1]
			intersections = append(intersections[:endIndex], intersections[endIndex+1:]...)
		}

		output = append(output, ObstacleSegment{
			Points: [2][2]float64{
				start,
				[2]float64{segments[i].Points[1][0], segments[i].Points[1][1]},
			},
			UserData: segments[i].UserData,
		})
	}
	return output
}
