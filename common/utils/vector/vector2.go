package vector

import (
	"encoding/json"
	"math"
	"math/rand"

	"github.com/bytearena/box2d"
	"github.com/go-gl/mathgl/mgl64"

	"github.com/bytearena/core/common/utils/number"
)

type Vector2 [2]float64
type Point2 = Vector2

func MakeVector2(x float64, y float64) Vector2 {
	return Vector2{x, y}
}

// Returns a random unit vector
func MakeRandomVector2() Vector2 {
	radians := rand.Float64() * math.Pi * 2
	return MakeVector2(
		math.Cos(radians),
		math.Sin(radians),
	)
}

// Returns a null vector2
func MakeNullVector2() Vector2 {
	return Vector2{0.0, 0.0}
}

func (v Vector2) Get() (float64, float64) {
	return v[0], v[1]
}

func (v Vector2) GetX() float64 {
	return v[0]
}

func (v Vector2) GetY() float64 {
	return v[1]
}

func (v Vector2) MarshalJSONString() string {
	json, _ := json.Marshal(v)
	return string(json)
}

func (a Vector2) Clone() Vector2 {
	return Vector2{a[0], a[1]}
}

func (a Vector2) Add(b Vector2) Vector2 {
	a[0] += b[0]
	a[1] += b[1]
	return a
}

func (a Vector2) AddScalar(f float64) Vector2 {
	if math.IsNaN(f) {
		return a
	}

	a[0] += f
	a[1] += f
	return a
}

func (a Vector2) Sub(b Vector2) Vector2 {
	a[0] -= b[0]
	a[1] -= b[1]
	return a
}

func (a Vector2) SubScalar(f float64) Vector2 {
	if math.IsNaN(f) {
		return a
	}

	a[0] -= f
	a[1] -= f
	return a
}

func (a Vector2) Scale(scale float64) Vector2 {
	if math.IsNaN(scale) {
		return a
	}

	a[0] *= scale
	a[1] *= scale
	return a
}

func (a Vector2) Mult(b Vector2) Vector2 {
	a[0] *= b[0]
	a[1] *= b[1]
	return a
}

func (a Vector2) MultScalar(f float64) Vector2 {
	if math.IsNaN(f) {
		return a
	}

	a[0] *= f
	a[1] *= f
	return a
}

func (a Vector2) Div(b Vector2) Vector2 {
	a[0] /= b[0]
	a[1] /= b[1]
	return a
}

func (a Vector2) DivScalar(f float64) Vector2 {
	if math.IsNaN(f) {
		return a
	}

	a[0] /= f
	a[1] /= f
	return a
}

func (a Vector2) Mag() float64 {
	return math.Sqrt(a.MagSq())
}

func (a Vector2) MagSq() float64 {
	return (a[0]*a[0] + a[1]*a[1])
}

func (a Vector2) SetMag(mag float64) Vector2 {
	if math.IsNaN(mag) {
		return a
	}

	return a.Normalize().MultScalar(mag)
}

func (a Vector2) Normalize() Vector2 {
	mag := a.Mag()
	if mag > 0 {
		return a.DivScalar(mag)
	}
	return a
}

func (a Vector2) OrthogonalClockwise() Vector2 {
	return MakeVector2(a[1], -a[0])
}

func (a Vector2) OrthogonalCounterClockwise() Vector2 {
	return MakeVector2(-a[1], a[0])
}

func (a Vector2) Center() Vector2 {
	return a.MultScalar(0.5)
}

func (a Vector2) Translate(translation Vector2) Vector2 {
	return a.Add(translation)
}

func (a Vector2) MoveCenterTo(newcenterpos Vector2) Vector2 {
	return a.Translate(newcenterpos.Sub(a.Center()))
}

func (a Vector2) SetAngle(radians float64) Vector2 {
	if math.IsNaN(radians) {
		return a
	}

	mag := a.Mag()
	a[0] = math.Sin(radians) * mag
	a[1] = math.Cos(radians) * mag

	if math.IsNaN(a[0]) || math.IsNaN(a[1]) {
		a[0], a[1] = 0, 0
	}

	return a
}

func (a Vector2) Limit(max float64) Vector2 {
	if math.IsNaN(max) {
		return a
	}

	mSq := a.MagSq()

	if mSq > max*max {
		return a.Normalize().MultScalar(max)
	}

	return a
}

func (a Vector2) Angle() float64 {
	if a[0] == 0 && a[1] == 0 {
		return 0
	}

	angle := math.Atan2(a[1], a[0])

	// Quart de tour à gauche
	angle = math.Pi/2.0 - angle

	if angle < 0 {
		angle += 2 * math.Pi
	}

	return angle
}

func (a Vector2) Cross(v Vector2) float64 {
	return a[0]*v[1] - a[1]*v[0]
}

func (a Vector2) Dot(v Vector2) float64 {
	return a[0]*v[0] + a[1]*v[1]
}

func (a Vector2) IsNull() bool {
	return isZero(a[0]) && isZero(a[1])
}

func (a Vector2) IsNullWithPrecision(precision float64) bool {
	return isZeroWithPrecision(a[0], precision) && isZeroWithPrecision(a[1], precision)
}

func (a Vector2) Equals(b Vector2) bool {
	return b.Sub(a).IsNull()
}

func (a Vector2) EqualsWithPrecision(b Vector2, precision float64) bool {
	return b.Sub(a).IsNullWithPrecision(precision)
}

func (a Vector2) String() string {
	return "<Vector2(" + number.FloatToStr(a[0], 5) + ", " + number.FloatToStr(a[1], 5) + ")>"
}

func (a Vector2) ToFloatArray() [2]float64 {
	return [2]float64{a.GetX(), a.GetY()}
}

func (a Vector2) ToB2Vec2() box2d.B2Vec2 {
	return box2d.MakeB2Vec2(a.GetX(), a.GetY())
}

func FromB2Vec2(v box2d.B2Vec2) Vector2 {
	return MakeVector2(v.X, v.Y)
}

var epsilon float64 = 0.000001

func isZero(f float64) bool {
	return math.Abs(f) < epsilon
}

func isZeroWithPrecision(f float64, precision float64) bool {
	return math.Abs(f) < precision
}

func (v Vector2) Transform(mat *mgl64.Mat4) Vector2 {
	res := mgl64.TransformCoordinate(
		mgl64.Vec3{
			v[0],
			0,
			v[1],
		},
		*mat,
	)
	return MakeVector2(res[0], res[2])
}
