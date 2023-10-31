//go:build example
// +build example

package main

import (
	"fmt"
	"image/color"

	imgui "github.com/gabstv/cimgui-go"
	"github.com/gabstv/ebiten-imgui/v2/renderer"
	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func main() {
	imctx := imgui.CreateContext()
	defer imctx.Destroy()

	gg := &G{
		c: renderer.NewCache(),
	}
	ebiten.SetWindowSize(800, 600)

	fonts := imgui.GetIO().GetFonts()
	_, _, _, _ = fonts.GetTextureDataAsRGBA32() // call this to force imgui to build the font atlas cache
	fonts.SetTexID(imgui.ImTextureID(1))

	ebiten.RunGame(gg)
}

type G struct {
	f float32
	c renderer.TextureCache
}

func (g *G) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{100, 100, 100, 255})
	io := imgui.GetIO()
	io.SetDisplaySize(imgui.ImVec2{X: 800, Y: 600})
	io.SetDeltaTime(1. / 60.)
	mx, my := ebiten.CursorPosition()
	io.SetMousePos(imgui.ImVec2{X: float32(mx), Y: float32(my)})
	io.SetMouseButtonDown(0, ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft))
	io.SetMouseButtonDown(1, ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight))
	io.SetMouseButtonDown(2, ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle))
	imgui.NewFrame()
	imgui.Text("Hello, world!")                // Display some text
	imgui.SliderFloat("float", &g.f, 0.0, 1.0) // Edit 1 float using a slider from 0.0f to 1.0f
	imgui.Render()

	renderer.Render(screen, imgui.GetDrawData(), g.c, ebiten.FilterNearest)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.2f", ebiten.CurrentTPS()))
}

func (g *G) Update() error {
	return nil
}

func (g *G) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}
