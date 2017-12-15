package vector

import (
	"fmt"
	"math"

	"github.com/bytearena/core/common/utils/number"
)

type Matrix2 [6]float64

func IdentityMatrix2() Matrix2 {
	return Matrix2{
		1, 0, 0,
		0, 1, 0,
	}
}

func MakeMatrix2() Matrix2 {
	return IdentityMatrix2()
}

func (m Matrix2) String() string {
	return fmt.Sprintf("<Matrix2[[%f,%f,%f],[%f,%f,%f]]>", m[0], m[1], m[2], m[3], m[4], m[5])
}

func (m Matrix2) IsIdentity() bool {
	id := IdentityMatrix2()
	return m[0] == id[0] &&
		m[1] == id[1] &&
		m[2] == id[2] &&
		m[3] == id[3] &&
		m[4] == id[4] &&
		m[5] == id[5]
}

func (m Matrix2) Translate(x, y float64) Matrix2 {
	return m.Mul(Matrix2{
		1, 0, x,
		0, 1, y,
	})
}

func (m Matrix2) Scale(x, y float64) Matrix2 {
	return m.Mul(Matrix2{
		x, 0, 0,
		0, y, 0,
	})
}

func (m Matrix2) Rotate(deg float64) Matrix2 {
	rad := number.DegreeToRadian(deg)
	sina := math.Sin(rad)
	cosa := math.Cos(rad)
	return m.Mul(Matrix2{
		cosa, -sina, 0,
		sina, cosa, 0,
	})
}

func (m Matrix2) Mul(n Matrix2) Matrix2 {
	return Matrix2{
		m[0]*n[0] + m[1]*n[3], m[0]*n[1] + m[1]*n[4], m[0]*n[2] + m[1]*n[5] + m[2],
		m[3]*n[0] + m[4]*n[3], m[3]*n[1] + m[4]*n[4], m[3]*n[2] + m[4]*n[5] + m[5],
	}
}

func (m Matrix2) Mulf(r float64) Matrix2 {
	for i := range m {
		m[i] *= r
	}
	return m
}

func (m Matrix2) Transform(x, y float64) (float64, float64) {
	return m[0]*x + m[1]*y + m[2], m[3]*x + m[4]*y + m[5]
}
