package browser

import (
	"bufio"
	"go-sandbox/gui/draw"
	"image/color"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.design/x/clipboard"
	"golang.org/x/net/html"
)

var (
	buttonX, buttonY          = 700, 20
	buttonWidth, buttonHeight = 100, 30
)

type Line struct {
	Text  string
	Depth int
}
type Browser struct {
	lines      []Line
	scrollY    int
	inputText  string
	lastSubmit string
	submitted  bool
	nodes      *html.Node
}

var _ ebiten.Game = (*Browser)(nil)

func NewBrowser() *Browser {
	return &Browser{lines: []Line{}}
}

func (b *Browser) submit() {
	b.submitted = true
	b.lastSubmit = b.inputText
	b.inputText = ""
	b.getHtml()
	walkDOM(b.nodes, 0, &b.lines)
}

func (b *Browser) Update() error {
	// scroll
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		b.scrollY += 20
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		b.scrollY -= 20
	}
	// mouse wheel
	_, dy := ebiten.Wheel()
	b.scrollY -= int(dy * 20)
	// スクロール範囲の制限
	if b.scrollY < 0 {
		b.scrollY = 0
	}

	// inputbox
	for _, r := range ebiten.AppendInputChars(nil) {
		b.inputText += string(r)
	}

	// バックスペース処理
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(b.inputText) > 0 {
		b.inputText = b.inputText[:len(b.inputText)-1]
	}

	// Enterで送信
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		if b.submitted || b.inputText == "" || !strings.Contains(b.inputText, "https://") {
			return nil
		}
		b.submit()
	}
	// ボタンを左クリックで送信
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if !b.submitted && b.inputText != "" && strings.Contains(b.inputText, "https://") {
			if buttonX <= x && x <= buttonX+buttonWidth && buttonY <= y && y <= buttonY+buttonHeight {
				b.submit()
			}
		}
	}

	isPaste := ebiten.IsKeyPressed(ebiten.KeyMeta) && inpututil.IsKeyJustPressed(ebiten.KeyV)
	if isPaste {
		pastedText := clipboard.Read(clipboard.FmtText)
		b.inputText += string(pastedText)
	}
	return nil
}

func (b *Browser) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	ttf, _ := os.ReadFile("./font_1_kokugl_1.15_rls.ttf")
	font, _ := truetype.Parse(ttf)
	face := truetype.NewFace(font, &truetype.Options{
		Size: 16,
	})

	if !b.submitted {
		// input box
		draw.InputBox(screen)
		text.Draw(screen, b.inputText, face, 30, 40, color.Black)

		// button
		clr := color.RGBA{0x80, 0x80, 0xff, 0xff}
		ebitenutil.DrawRect(screen, float64(buttonX), float64(buttonY), float64(buttonWidth), float64(buttonHeight), clr)

	}

	if b.submitted {
		y := 30
		for _, line := range b.lines {
			// x := 20 + line.Depth*20
			x := 20
			text.Draw(screen, line.Text, face, x, y-b.scrollY, color.Black)
			y += 20
		}
	}
}

func (b *Browser) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (b *Browser) getHtml() error {
	data, _ := http.Get(b.lastSubmit)
	defer data.Body.Close()

	doc, err := html.Parse(bufio.NewReader(data.Body))
	if err != nil {
		log.Fatal(err)
		return err
	}

	b.nodes = doc
	return nil
}
