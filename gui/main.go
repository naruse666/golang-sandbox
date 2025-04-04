package main

import (
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/net/html"
)

type Game struct {
	texts []string
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	face := basicfont.Face7x13
	y := 30
	for _, line := range g.texts {
		text.Draw(screen, line, face, 20, y, color.White)
		y += 20
	}
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}

func main() {
	doc, err := html.Parse(strings.NewReader("<html><body><p>hello</p><div>Div</div></body></html>"))

	if err != nil {
		log.Fatal(err)
	}

	var texts []string
	extractPTexts(doc, &texts)

	game := &Game{texts: texts}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("test window")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func extractPTexts(n *html.Node, results *[]string) {
	if n.Type == html.ElementNode {
		if n.Data == "p" || n.Data == "div" {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				if c.Type == html.TextNode {
					*results = append(*results, c.Data)
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractPTexts(c, results)
	}
}
