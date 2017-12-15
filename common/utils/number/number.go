package number

import (
	"math"
	"strconv"
	"time"
)

func ToFixed(val float64, places int) (newVal float64) {
	roundOn := 0.5
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	if newVal == -0.0 {
		newVal = 0.0
	}
	return
}

func FloatToStr(f float64, places int) string {
	return strconv.FormatFloat(f, 'f', places, 64)
}

func DiffMs(b time.Time, a time.Time) float64 {
	return float64(b.UnixNano()-a.UnixNano()) / 1000000.0
}

func DurationMs(d time.Duration) float64 {
	return float64(d.Nanoseconds()) / 1000000.0
}

func Map(n float64, start1 float64, stop1 float64, start2 float64, stop2 float64) float64 {
	return ((n-start1)/(stop1-start1))*(stop2-start2) + start2
}

func RadianToDegree(rad float64) float64 {
	return rad * (180.0 / math.Pi)
}

func DegreeToRadian(deg float64) float64 {
	return deg * (math.Pi / 180.0)
}

var epsilon float64 = 0.0000000001

func IsZero(f float64) bool {
	return math.Abs(f) < epsilon
}

func FloatEquals(a, b float64) bool {
	return IsZero(b - a)
}

func Round(val float64) int {
	if val < 0 {
		return int(val - 0.5)
	}
	return int(val + 0.5)
}
