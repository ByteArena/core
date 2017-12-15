package trigo

import (
	"errors"

	"github.com/dhconnelly/rtreego"

	"github.com/bytearena/core/common/utils/vector"
)

func makePoly(centerA, centerB vector.Vector2, radiusA, radiusB float64) []vector.Vector2 {
	hasMoved := !centerA.Equals(centerB)

	if !hasMoved {
		return nil
	}

	AB := vector.MakeSegment2(centerA, centerB)

	// on détermine les 4 points formant le rectangle orienté définissant la trajectoire de l'object en movement

	polyASide := AB.OrthogonalToACentered().SetLengthFromCenter(radiusA * 2) // si AB vertical, A à gauche, B à droite
	polyBSide := AB.OrthogonalToBCentered().SetLengthFromCenter(radiusB * 2) // si AB vertical, A à gauche, B à droite

	/*

		B2*--------------------*A2
		  |                    |
		B *                    * A
		  |                    |
		B1*--------------------*A1

	*/

	polyA1, polyA2 := polyASide.Get()
	polyB1, polyB2 := polyBSide.Get()

	return []vector.Vector2{polyA1, polyA2, polyB2, polyB1}
}

func getGeometryObjectBoundingBox(position vector.Vector2, radius float64) (bottomLeft vector.Vector2, topRight vector.Vector2) {
	x, y := position.Get()
	return vector.MakeVector2(x-radius, y-radius), vector.MakeVector2(x+radius, y+radius)
}

func GetTrajectoryBoundingBox(beginPoint vector.Vector2, beginRadius float64, endPoint vector.Vector2, endRadius float64) (*rtreego.Rect, error) {
	beginBottomLeft, beginTopRight := getGeometryObjectBoundingBox(beginPoint, beginRadius)
	endBottomLeft, endTopRight := getGeometryObjectBoundingBox(endPoint, endRadius)

	bbTopLeft, bbDimensions := GetBoundingBox([]vector.Vector2{beginBottomLeft, beginTopRight, endBottomLeft, endTopRight})

	//show := spew.ConfigState{MaxDepth: 5, Indent: "    "}

	bbRegion, err := rtreego.NewRect(bbTopLeft, bbDimensions)
	if err != nil {
		return nil, errors.New("Error in getTrajectoryBoundingBox: could not define bbRegion in rTree")
	}

	// fmt.Println("----------------------------------------------------------------")
	// show.Dump(bbTopLeft, bbDimensions, bbRegion)
	// fmt.Println("----------------------------------------------------------------")

	return bbRegion, nil
}

func GetBoundingBox(points []vector.Vector2) (rtreego.Point, rtreego.Point) {

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

	return []float64{minX, minY}, []float64{width, height}
}
