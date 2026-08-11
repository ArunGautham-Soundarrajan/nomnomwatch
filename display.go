package main

import "image/color"

type Displayer interface {
	SetPixel(x, y int16, c color.RGBA)
	FillScreen(c color.Color)
}

func Render(d Displayer) {
	d.FillScreen(color.White)
}
