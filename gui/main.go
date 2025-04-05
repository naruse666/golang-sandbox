package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	// "net/http"
	"github.com/golang/freetype/truetype"
	"os"
	"strings"

	"github.com/kbinani/screenshot"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"golang.org/x/net/html"
)

type Game struct {
	lines   []Line
	scrollY int
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.scrollY += 20
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.scrollY -= 20
	}
	// mouse wheel
	_, dy := ebiten.Wheel()
	g.scrollY -= int(dy * 20)
	// スクロール範囲の制限
	if g.scrollY < 0 {
		g.scrollY = 0
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ttf, _ := os.ReadFile("./font_1_kokugl_1.15_rls.ttf")
	font, _ := truetype.Parse(ttf)
	face := truetype.NewFace(font, &truetype.Options{
		Size: 16,
	})
	y := 30
	for _, line := range g.lines {
		// x := 20 + line.Depth*20
		x := 20
		text.Draw(screen, line.Text, face, x, y-g.scrollY, color.White)
		y += 20
	}
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	data, err := os.Open("test.html")
	if err != nil {
		panic(err)
	}
	doc, err := html.Parse(bufio.NewReader(data))
	// doc, err := html.Parse(strings.NewReader("<html><body><p>hello</p><div>Div</div></body></html>"))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v", doc)
	var lines []Line
	walkDOM(doc, 0, &lines)

	game := &Game{lines: lines}
	width, height := getPrimaryScreenSize()
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("test window")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

type Line struct {
	Text  string
	Depth int
}

func walkDOM(n *html.Node, depth int, lines *[]Line) {
	switch n.Type {
	case html.ElementNode:
		// *lines = append(*lines, Line{Text: n.Data, Depth: depth})
		for _, a := range n.Attr {
			if a.Key == "class" {
				fmt.Println(n.Data, a.Key, a.Val)
			}
		}
	case html.TextNode:
		trimmed := strings.TrimSpace(n.Data)
		if trimmed != "" {
			*lines = append(*lines, Line{Text: trimmed, Depth: depth})
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkDOM(c, depth+1, lines)
	}
}

func getPrimaryScreenSize() (int, int) {
	n := screenshot.NumActiveDisplays()
	if n <= 0 {
		log.Fatal("no display detected")
	}

	// ここでは一番左上のディスプレイを取得
	bounds := screenshot.GetDisplayBounds(0)
	width := bounds.Dx()
	height := bounds.Dy()
	return width, height
}
