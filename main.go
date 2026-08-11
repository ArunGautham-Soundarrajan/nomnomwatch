package main

import (
	"time"

	"github.com/sago35/tinydisplay"
)

// 240x135 (LilyGo T-Display S3) or 240x240, 320x240, etc.
const (
	displayWidth  = 240
	displayHeight = 135
)

func main() {
	display, _ := tinydisplay.NewClient("127.0.0.1", 9812, displayWidth, displayHeight)

	Render(display)
	time.Sleep(time.Minute)
}
