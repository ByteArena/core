package polygon

import (
	"errors"

	"github.com/bytearena/core/common/utils/vector"
)

var CartesianSystemWinding = (struct {
	CW  int
	CCW int
}{
	CW:  1,
	CCW: -1,
})

func GetPolygonWindingForCartesianSystem(poly []vector.Vector2) int {
	polylen := len(poly)
	if polylen < 3 {
		return 0
	}

	sum := 0.0

	for i := 0; i < polylen-1; i++ { // -1: last point is first point

		p1 := poly[i]
		p2 := poly[(i+1)%polylen]

		sum += (p2.GetX() - p1.GetX()) * (p2.GetY() + p1.GetY())
	}

	if sum > 0 {
		// CW
		return CartesianSystemWinding.CW
	}

	if sum < 0 {
		// CCW
		return CartesianSystemWinding.CCW
	}

	return 0
}

func IsCW(winding int) bool {
	return winding == CartesianSystemWinding.CW
}

func IsCCW(winding int) bool {
	return winding == CartesianSystemWinding.CCW
}

func InvertWinding(poly []vector.Vector2) []vector.Vector2 {
	new := make([]vector.Vector2, len(poly))

	newindex := 0
	for oldindex := len(poly) - 1; oldindex >= 0; oldindex-- {
		point := poly[oldindex]
		new[newindex] = point
		newindex++
	}

	return new
}

func EnsureWinding(winding int, poly []vector.Vector2) ([]vector.Vector2, error) {

	curwinding := GetPolygonWindingForCartesianSystem(poly)

	// Change polygons winding to CW
	if curwinding != winding {
		// Change winding
		newpoly := InvertWinding(poly)
		newwinding := GetPolygonWindingForCartesianSystem(newpoly)
		if newwinding != winding {
			return nil, errors.New("Could not change polygon winding from CW to CCW")
		}

		return newpoly, nil
	}

	return poly, nil
}
