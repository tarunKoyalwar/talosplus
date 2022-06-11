package ioutils

import (
	"image/color"

	"github.com/muesli/termenv"
)

var (
	Green      = GetColor(126, 211, 33)
	Orange     = GetColor(245, 166, 35)
	AquaMarine = GetColor(80, 227, 194)
	Azure      = GetColor(74, 144, 226)
	Yellow     = GetColor(248, 231, 28)
	Grey       = GetColor(155, 155, 155)
	Red        = GetColor(208, 2, 27)
	LightGreen = GetColor(184, 233, 134)
)

var Profile termenv.Profile = termenv.ColorProfile()

func GetColor(r uint8, g uint8, b uint8) termenv.Color {
	return Profile.FromColor(color.RGBA{R: r, G: g, B: b, A: 255})

}
