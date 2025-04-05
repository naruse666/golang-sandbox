package main

import (
	"bufio"
	"go-sandbox/gui/draw"
	"image/color"
	"log"
	"net/http"

	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kbinani/screenshot"
	"golang.design/x/clipboard"

	"golang.org/x/net/html"
)

type Content struct {
	lines      []Line
	scrollY    int
	inputText  string
	lastSubmit string
	submitted  bool
	nodes      *html.Node
}

var (
	buttonX, buttonY          = 700, 20
	buttonWidth, buttonHeight = 100, 30
)

func (c *Content) submit() {
	c.submitted = true
	c.lastSubmit = c.inputText
	c.inputText = ""
	c.getHtml()
	walkDOM(c.nodes, 0, &c.lines)
}

func (c *Content) Update() error {
	// scroll
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		c.scrollY += 20
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		c.scrollY -= 20
	}
	// mouse wheel
	_, dy := ebiten.Wheel()
	c.scrollY -= int(dy * 20)
	// スクロール範囲の制限
	if c.scrollY < 0 {
		c.scrollY = 0
	}

	// inputbox
	for _, r := range ebiten.AppendInputChars(nil) {
		c.inputText += string(r)
	}

	// バックスペース処理
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(c.inputText) > 0 {
		c.inputText = c.inputText[:len(c.inputText)-1]
	}

	// Enterで送信
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if c.submitted || c.inputText == "" || !strings.Contains(c.inputText, "https://") {
			return nil
		}
		c.submit()
	}
	// ボタンを左クリックで送信
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if !c.submitted && c.inputText != "" && strings.Contains(c.inputText, "https://") {
			if buttonX <= x && x <= buttonX+buttonWidth && buttonY <= y && y <= buttonY+buttonHeight {
				c.submit()
			}
		}
	}

	isPaste := ebiten.IsKeyPressed(ebiten.KeyMeta) && inpututil.IsKeyJustPressed(ebiten.KeyV)
	if isPaste {
		pastedText := clipboard.Read(clipboard.FmtText)
		c.inputText += string(pastedText)
	}
	return nil
}

func (c *Content) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	ttf, _ := os.ReadFile("./font_1_kokugl_1.15_rls.ttf")
	font, _ := truetype.Parse(ttf)
	face := truetype.NewFace(font, &truetype.Options{
		Size: 16,
	})

	if !c.submitted {
		// input box
		draw.InputBox(screen)
		text.Draw(screen, c.inputText, face, 30, 40, color.Black)

		// button
		clr := color.RGBA{0x80, 0x80, 0xff, 0xff}
		ebitenutil.DrawRect(screen, float64(buttonX), float64(buttonY), float64(buttonWidth), float64(buttonHeight), clr)

	}

	if c.submitted {
		y := 30
		for _, line := range c.lines {
			// x := 20 + line.Depth*20
			x := 20
			text.Draw(screen, line.Text, face, x, y-c.scrollY, color.Black)
			y += 20
		}
	}
}
func (c *Content) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (c *Content) getHtml() error {
	data, _ := http.Get(c.lastSubmit)
	defer data.Body.Close()

	doc, err := html.Parse(bufio.NewReader(data.Body))
	if err != nil {
		log.Fatal(err)
		return err
	}

	c.nodes = doc
	return nil
}

func main() {
	// data, err := os.Open("test.html")
	// if err != nil {
	// 	panic(err)
	// }
	// doc, err := html.Parse(strings.NewReader("<html><body><p>hello</p><div>Div</div></body></html>"))

	var lines []Line

	content := &Content{lines: lines}
	width, height := getPrimaryScreenSize()
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("test window")

	// ebiten.SetMaxTPS(10)
	if err := ebiten.RunGame(content); err != nil {
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
		// for _, a := range n.Attr {
		// 	if a.Key == "class" {
		// 		fmt.Println(n.Data, a.Key, a.Val)
		// 	}
		// }
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
