//go:build example
// +build example

package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	imgui "github.com/gabstv/cimgui-go"
	"github.com/gabstv/ebiten-imgui/v2/renderer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var exampleImage *ebiten.Image

func main() {
	mgr := renderer.New(nil)

	ebiten.SetWindowSize(800, 600)

	gg := &G{
		mgr: mgr,
	}

	exampleImage, _, err := ebitenutil.NewImageFromFile("example.png")
	if err != nil {
		log.Fatal(err)
	}
	mgr.Cache.SetTexture(10, exampleImage) // Texture ID 10 will contain this example image

	ebiten.RunGame(gg)
}

type G struct {
	mgr *renderer.Manager
}

func (g *G) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.2f", ebiten.CurrentTPS()))
	g.mgr.Draw(screen)
}

// cimgui-go doesn't have imgui.Image with sensible default yet
// you can implement it like this
func Image(tid imgui.ImTextureID, size imgui.ImVec2) {
	uv0 := imgui.NewImVec2(0, 0)
	uv1 := imgui.NewImVec2(1, 1)
	border_col := imgui.NewImVec4(0, 0, 0, 0)
	tint_col := imgui.NewImVec4(1, 1, 1, 1)

	imgui.ImageV(tid, size, uv0, uv1, tint_col, border_col)
}

func (g *G) Update() error {
	g.mgr.Update(1.0 / 60.0)
	g.mgr.BeginFrame()
	imgui.Text("Hello, images!")
	Image(10, imgui.ImVec2{X: 64, Y: 64})
	g.mgr.EndFrame()
	return nil
}

func (g *G) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.mgr.SetDisplaySize(float32(800), float32(600))
	return 800, 600
}
