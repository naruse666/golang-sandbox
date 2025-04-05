package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/kbinani/screenshot"
	"go-sandbox/gui/browser"
)

func main() {

	browser := browser.NewBrowser()
	width, height := getPrimaryScreenSize()
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("test window")

	// ebiten.SetMaxTPS(10)
	if err := ebiten.RunGame(browser); err != nil {
		panic(err)
	}
}

func getPrimaryScreenSize() (int, int) {
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		log.Fatal("no display detected")
	}

	bounds := screenshot.GetDisplayBounds(0)
	width := bounds.Dx()
	height := bounds.Dy()
	return width, height
}
