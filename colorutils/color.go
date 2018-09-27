package colorutils

import (
	"github.com/lucasb-eyer/go-colorful"
	"image/color"
)

func IsBlackOrWhite(color color.Color) bool {
	colorNew, _ := colorful.MakeColor(color)

	if (colorNew.R == 0 && colorNew.G == 0 && colorNew.B == 0) || (colorNew.R == 1 && colorNew.G == 1 && colorNew.B == 1) {
		return true
	}
	return false
}

func IsSameColor(color1 color.Color, color2 color.Color) bool {

	color1New, _ := colorful.MakeColor(color1)
	color2New, _ := colorful.MakeColor(color2)

	if color1New.R == color2New.R && color1New.G == color2New.G && color1New.B == color2New.B {
		return true
	}
	return false
}