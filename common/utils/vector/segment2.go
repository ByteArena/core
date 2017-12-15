package vector

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Segment2 [2]Vector2

func MakeSegment2(a Vector2, b Vector2) Segment2 {
	return Segment2{
		a,
		b,
	}
}

func (s Segment2) Get() (Vector2, Vector2) {
	return s[0], s[1]
}

func (s Segment2) GetPointA() Vector2 {
	return s[0]
}

func (s Segment2) GetPointB() Vector2 {
	return s[1]
}

func (s Segment2) SetPointA(v Vector2) Segment2 {
	s[0] = v
	return s
}

func (s Segment2) SetPointB(v Vector2) Segment2 {
	s[1] = v
	return s
}

func (s Segment2) Equals(s2 Segment2) bool {
	return s2[0].Equals(s[0]) && s2[1].Equals(s[1])
}

func (s Segment2) String() string {
	return "<Segment2(" + s[0].MarshalJSONString() + ", " + s[1].MarshalJSONString() + ")>"
}

func (s Segment2) Clone() Segment2 {
	return MakeSegment2(s[0].Clone(), s[1].Clone())
}

func (s Segment2) Add(v Vector2) Segment2 {
	s[0] = s[0].Add(v)
	s[1] = s[1].Add(v)
	return s
}

func (s Segment2) AddScalar(f float64) Segment2 {
	s[0] = s[0].AddScalar(f)
	s[1] = s[1].AddScalar(f)
	return s
}

func (s Segment2) Sub(v Vector2) Segment2 {
	s[0] = s[0].Sub(v)
	s[1] = s[1].Sub(v)
	return s
}

func (s Segment2) SubScalar(f float64) Segment2 {
	s[0] = s[0].SubScalar(f)
	s[1] = s[1].SubScalar(f)
	return s
}

// Vector2 returns Vector2 from a to b relative to (0,0)
func (s Segment2) Vector2() Vector2 {
	return s[1].Sub(s[0])
}

func (s Segment2) Center() Vector2 {
	relativeCenter := s.Vector2().MultScalar(0.5)
	return s[0].Add(relativeCenter)
}

func (s Segment2) Translate(v Vector2) Segment2 {
	s[0] = s[0].Add(v)
	s[1] = s[1].Add(v)
	return s
}

func (s Segment2) ScaleFromA(scale float64) Segment2 {
	s[1] = s[1].Sub(s[0]).Scale(scale).Add(s[0])
	return s
}

func (s Segment2) ScaleFromB(scale float64) Segment2 {
	s[0] = s[0].Sub(s[1]).Scale(scale).Add(s[1])
	return s
}

func (s Segment2) ScaleFromCenter(scale float64) Segment2 {
	center := s.Center()
	s[1] = s[1].Sub(center).Scale(scale).Add(center)
	s[0] = s[0].Sub(center).Scale(scale).Add(center)
	return s
}

func (s Segment2) LengthSq() float64 {
	return s.Vector2().MagSq()
}

func (s Segment2) Length() float64 {
	return s.Vector2().Mag()
}

func (s Segment2) NormalizeFromA() Segment2 {
	normalized := s.Vector2().Normalize()
	s[1] = s[0].Add(normalized)
	return s
}

func (s Segment2) NormalizeFromB() Segment2 {
	normalized := s.Vector2().Normalize()
	s[0] = s[1].Sub(normalized)
	return s
}

func (s Segment2) NormalizeFromCenter() Segment2 {
	halfnormalized := s.Vector2().Normalize().Scale(0.5)
	center := s.Center()
	s[0] = center.Add(halfnormalized)
	s[1] = center.Sub(halfnormalized)
	return s
}

func (s Segment2) SetLengthFromA(length float64) Segment2 {
	return s.NormalizeFromA().ScaleFromA(length)
}

func (s Segment2) SetLengthFromB(length float64) Segment2 {
	return s.NormalizeFromB().ScaleFromB(length)
}

func (s Segment2) SetLengthFromCenter(length float64) Segment2 {
	return s.NormalizeFromCenter().ScaleFromCenter(length)
}

func (s Segment2) OrthogonalToAClockwise() Segment2 {
	ortho := s.Vector2().OrthogonalClockwise()
	s[1] = s[0].Add(ortho)
	return s
}

func (s Segment2) OrthogonalToACounterClockwise() Segment2 {
	ortho := s.Vector2().OrthogonalCounterClockwise()
	s[1] = s[0].Add(ortho)
	return s
}

func (s Segment2) OrthogonalToACentered() Segment2 {
	newS := s.OrthogonalToAClockwise()
	return newS.Translate(newS.Vector2().Scale(-0.5))
}

func (s Segment2) OrthogonalToBClockwise() Segment2 {
	ortho := s.Vector2().OrthogonalClockwise()
	s[0] = s[1]
	s[1] = s[1].Add(ortho)
	return s
}

func (s Segment2) OrthogonalToBCounterClockwise() Segment2 {
	ortho := s.Vector2().OrthogonalCounterClockwise()
	s[0] = s[1]
	s[1] = s[1].Add(ortho)
	return s
}

func (s Segment2) OrthogonalToBCentered() Segment2 {
	newS := s.OrthogonalToBClockwise()
	return newS.Translate(newS.Vector2().Scale(-0.5))
}

func (s Segment2) OrthogonalToCenterClockwise() Segment2 {
	ortho := s.Vector2().OrthogonalClockwise()
	center := s.Center()
	s[0] = center
	s[1] = s[0].Add(ortho)
	return s
}

func (s Segment2) OrthogonalToCenterCounterClockwise() Segment2 {
	ortho := s.Vector2().OrthogonalCounterClockwise()
	center := s.Center()
	s[0] = center
	s[1] = s[0].Add(ortho)
	return s
}

func (s Segment2) OrthogonalToCenterCentered() Segment2 {
	newS := s.OrthogonalToCenterClockwise()
	return newS.Translate(newS.Vector2().Scale(-0.5))
}

func (s Segment2) MoveCenterTo(newcenterpos Vector2) Segment2 {
	center := s.Center()
	translation := newcenterpos.Sub(center)
	return s.Translate(translation)
}

func (s Segment2) ToRectangleCentered(height float64) []Vector2 {
	a1, a2 := s.OrthogonalToACentered().SetLengthFromCenter(height).Get()
	b1, b2 := s.OrthogonalToBCentered().SetLengthFromCenter(height).Get()

	return []Vector2{
		a1,
		b1,
		b2,
		a2,
	}
}

var testnum int

func test(ok bool, testname string) {
	testnum++
	if !ok {
		fmt.Println("FAILED #" + strconv.Itoa(testnum) + ": " + testname)
		os.Exit(1)
	}

	fmt.Println("SUCCESS #" + strconv.Itoa(testnum) + ": " + testname)
}

func TestSegment2() {
	log.Println("Testing Segment2")

	va := MakeVector2(-1.5, 3.5)
	vb := MakeVector2(-3, 2.5)

	var sExpected Segment2
	var vExpected Vector2
	s := MakeSegment2(va, vb)

	// Clone
	sclone := s.Clone()

	test(s[0].Equals(sclone[0]), "Cloned [0]")
	test(s[1].Equals(sclone[1]), "Cloned [1]")

	// Equals
	test(sclone.Equals(s), "Equals")

	// JSON
	json, _ := json.Marshal(s)
	test(string(json) == "[[-1.5000,3.5000],[-3.0000,2.5000]]", "JSON")

	// Add
	vadd := MakeVector2(-1, 10)
	sExpected = MakeSegment2(
		va.Add(vadd),
		vb.Add(vadd),
	)
	test(s.Add(vadd).Equals(sExpected), "Add")

	// AddScalar
	sExpected = MakeSegment2(
		va.AddScalar(8.444),
		vb.AddScalar(8.444),
	)
	test(s.AddScalar(8.444).Equals(sExpected), "AddScalar")

	// Sub
	sExpected = MakeSegment2(
		va.Sub(vadd),
		vb.Sub(vadd),
	)
	test(s.Sub(vadd).Equals(sExpected), "Sub")

	// SubScalar
	sExpected = MakeSegment2(
		va.SubScalar(8.444),
		vb.SubScalar(8.444),
	)
	test(s.SubScalar(8.444).Equals(sExpected), "SubScalar")

	// Vector2
	vExpected = vb.Sub(va)
	test(s.Vector2().Equals(vExpected), "Vector2")

	// Center
	vExpected = MakeVector2(-2.25, 3)
	test(s.Center().Equals(vExpected), "Center")

	// Translate
	vtranslate := MakeVector2(-123, 10)
	sExpected = MakeSegment2(
		MakeVector2(-124.5, 13.5),
		MakeVector2(-126, 12.5),
	)
	test(s.Translate(vtranslate).Equals(sExpected), "Translate")

	// ScaleFromA
	sExpected = MakeSegment2(
		va,
		MakeVector2(-4.5, 1.5),
	)
	test(s.ScaleFromA(2).Equals(sExpected), "ScaleFromA")

	// ScaleFromB
	sExpected = MakeSegment2(
		MakeVector2(0, 4.5),
		vb,
	)
	test(s.ScaleFromB(2).Equals(sExpected), "ScaleFromB")

	// ScaleFromCenter
	sExpected = MakeSegment2(
		MakeVector2(-0.75, 4),
		MakeVector2(-3.75, 2),
	)
	test(s.ScaleFromCenter(2).Equals(sExpected), "ScaleFromCenter")

	// LengthSq
	fexpected := 1.8027756377319946 * 1.8027756377319946
	test(isZero(s.LengthSq()-fexpected), "LengthSq")

	// LengthSq
	fexpected = 1.8027756377319946
	test(isZero(s.Length()-fexpected), "Length")

	// NormalizeFromA
	normalized := s.Vector2().Normalize()
	sExpected = s
	sExpected[1] = sExpected[0].Add(normalized)
	test(isZero(s.NormalizeFromA().Length()-1.0), "NormalizeFromA:Length")
	test(s.NormalizeFromA().Equals(sExpected), "NormalizeFromA")

	// NormalizeFromB
	sExpected = s
	sExpected[0] = sExpected[1].Sub(normalized)
	test(isZero(s.NormalizeFromB().Length()-1.0), "NormalizeFromB:Length")
	test(s.NormalizeFromB().Equals(sExpected), "NormalizeFromB")

	// NormalizeFromCenter
	sExpected = MakeSegment2(
		MakeVector2(-2.666025, 2.72265),
		MakeVector2(-1.833975, 3.27735),
	)
	test(isZero(s.NormalizeFromCenter().Length()-1.0), "NormalizeFromCenter:Length")
	test(s.NormalizeFromCenter().Equals(sExpected), "NormalizeFromCenter")

	// SetLengthFromA
	sExpected = MakeSegment2(
		MakeVector2(-1.5, 3.5),
		MakeVector2(-1.5+-0.416025, 3.5+-0.27735),
	)
	test(isZero(s.SetLengthFromA(0.5).Length()-0.5), "SetLengthFromA:Length")
	test(s.SetLengthFromA(0.5).Equals(sExpected), "SetLengthFromA")

	// SetLengthFromB
	sExpected = MakeSegment2(
		MakeVector2(-3.0 - -0.416025, 2.5 - -0.27735),
		MakeVector2(-3.0, 2.5),
	)
	test(isZero(s.SetLengthFromB(0.5).Length()-0.5), "SetLengthFromB:Length")
	test(s.SetLengthFromB(0.5).Equals(sExpected), "SetLengthFromB")

	// SetLengthFromCenter
	sExpected = s.NormalizeFromCenter()
	test(isZero(s.SetLengthFromCenter(1).Length()-1.0), "SetLengthFromCenter:Length")
	test(s.SetLengthFromCenter(1).Equals(sExpected), "SetLengthFromCenter")

	// OrthogonalToAClockwise
	sExpected = MakeSegment2(
		MakeVector2(-1.5000, 3.5000),
		MakeVector2(-2.5000, 5.0000),
	)
	test(s.OrthogonalToAClockwise().Equals(sExpected), "OrthogonalToAClockwise")

	// OrthogonalToACounterClockwise
	sExpected = MakeSegment2(
		MakeVector2(-1.5000, 3.5000),
		MakeVector2(-0.5000, 2.0000),
	)
	test(s.OrthogonalToACounterClockwise().Equals(sExpected), "OrthogonalToACounterClockwise")

	// OrthogonalToACentered
	sExpected = MakeSegment2(
		MakeVector2(-1.0000, 2.7500),
		MakeVector2(-2.0000, 4.2500),
	)
	test(s.OrthogonalToACentered().Equals(sExpected), "OrthogonalToACentered")

	// OrthogonalToBClockwise
	sExpected = MakeSegment2(
		MakeVector2(-3, 2.5),
		MakeVector2(-4, 4),
	)
	test(s.OrthogonalToBClockwise().Equals(sExpected), "OrthogonalToBClockwise")

	// OrthogonalToBCounterClockwise
	sExpected = MakeSegment2(
		MakeVector2(-3, 2.5),
		MakeVector2(-2, 1),
	)
	test(s.OrthogonalToBCounterClockwise().Equals(sExpected), "OrthogonalToBCounterClockwise")

	// OrthogonalToBCentered
	sExpected = MakeSegment2(
		MakeVector2(-2.5000, 1.7500),
		MakeVector2(-3.5000, 3.2500),
	)
	test(s.OrthogonalToBCentered().Equals(sExpected), "OrthogonalToBCentered")

	// OrthogonalToCenterClockwise
	sExpected = MakeSegment2(
		MakeVector2(-2.2500, 3.0000),
		MakeVector2(-3.2500, 4.5000),
	)
	test(s.OrthogonalToCenterClockwise().Equals(sExpected), "OrthogonalToCenterClockwise")

	// OrthogonalToCenterCounterClockwise
	sExpected = MakeSegment2(
		MakeVector2(-2.2500, 3.0000),
		MakeVector2(-1.2500, 1.5000),
	)
	test(s.OrthogonalToCenterCounterClockwise().Equals(sExpected), "OrthogonalToCenterCounterClockwise")

	// OrthogonalToCenterCentered
	sExpected = MakeSegment2(
		MakeVector2(-1.7500, 2.2500),
		MakeVector2(-2.7500, 3.7500),
	)
	test(s.OrthogonalToCenterCentered().Equals(sExpected), "OrthogonalToCenterCentered")

	// MoveCenterTo
	sExpected = MakeSegment2(
		MakeVector2(0.75, 0.5),
		MakeVector2(-0.75, -0.5),
	)
	test(s.MoveCenterTo(MakeVector2(0, 0)).Equals(sExpected), "MoveCenterTo")
}
