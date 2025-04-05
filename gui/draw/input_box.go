package draw

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

func InputBox(screen *ebiten.Image) {
	boxColor := color.RGBA{230, 230, 230, 255}
	inputBox := ebiten.NewImage(600, 30)
	inputBox.Fill(boxColor)
	inputBox.Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(20, 20)
	screen.DrawImage(inputBox, op)
}
