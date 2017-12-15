package types

import "github.com/bytearena/core/common/utils/number"

type PixelUnit float64

func (p PixelUnit) Pixels() float64 {
	return float64(p)
}

func (p PixelUnit) RoundPixels() int {
	return number.Round(p.Pixels())
}

type PixelSurface struct {
	Width  PixelUnit
	Height PixelUnit
}
