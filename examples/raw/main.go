// +build example

package main

import (
	"fmt"
	"image/color"
	"unsafe"

	"github.com/gabstv/ebiten-imgui/renderer"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/inkyblackness/imgui-go/v2"
)

func main() {
	imctx := imgui.CreateContext(nil)
	defer imctx.Destroy()
	io := imgui.CurrentIO()
	//io.SetClipboard()
	io.SetClipboard(clipboard{})
	gg := &G{
		c: make(map[imgui.TextureID]*ebiten.Image),
	}
	ebiten.SetWindowSize(800, 600)

	// Build texture atlas
	image := io.Fonts().TextureDataAlpha8()
	rawimg := make([]uint8, image.Width*image.Height)
	for i := range rawimg {
		rawimg[i] = 255
	}
	image.Pixels = unsafe.Pointer(&rawimg[0])
	io.Fonts().SetTextureID(1)

	ebiten.RunGame(gg)
}

type G struct {
	f float32
	c map[imgui.TextureID]*ebiten.Image
}

func (g *G) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{100, 100, 100, 255})
	io := imgui.CurrentIO()
	io.SetDisplaySize(imgui.Vec2{X: 800, Y: 600})
	io.SetDeltaTime(1. / 60.)
	mx, my := ebiten.CursorPosition()
	io.SetMousePosition(imgui.Vec2{X: float32(mx), Y: float32(my)})
	io.SetMouseButtonDown(0, ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft))
	io.SetMouseButtonDown(1, ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight))
	io.SetMouseButtonDown(2, ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle))
	imgui.NewFrame()
	imgui.Text("Hello, world!")                // Display some text
	imgui.SliderFloat("float", &g.f, 0.0, 1.0) // Edit 1 float using a slider from 0.0f to 1.0f
	imgui.Render()

	renderer.Render(screen, imgui.RenderedDrawData(), g.c)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.2f", ebiten.CurrentTPS()))
}

func (g *G) Update(screen *ebiten.Image) error {
	return nil
}

func (g *G) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}

type clipboard struct {
	//platform Platform
}

func (board clipboard) Text() (string, error) {
	return "", nil //board.platform.ClipboardText()
}

func (board clipboard) SetText(text string) {
	//board.platform.SetClipboardText(text)
}
