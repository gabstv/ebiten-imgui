// +build example

package main

import (
	"fmt"
	"image/color"

	"github.com/gabstv/ebiten-imgui/renderer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/inkyblackness/imgui-go/v2"
)

var exampleImage *ebiten.Image

func main() {
	mgr := renderer.New(nil)

	ebiten.SetWindowSize(800, 600)

	gg := &G{
		mgr: mgr,
	}

	exampleImage, _, _ = ebitenutil.NewImageFromFile("example.png", ebiten.FilterNearest)
	mgr.Cache.SetTexture(10, exampleImage) // Texture ID 10 will contain this example image

	ebiten.RunGame(gg)
}

type G struct {
	mgr *renderer.Manager
	// demo members:
	clearColor [3]float32
}

func (g *G) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{uint8(g.clearColor[0] * 255), uint8(g.clearColor[1] * 255), uint8(g.clearColor[2] * 255), 255})
	g.mgr.BeginFrame()

	{
		imgui.Text("Hello, images!")
		imgui.Image(10, imgui.Vec2{64, 64})

	}

	g.mgr.EndFrame(screen)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.2f", ebiten.CurrentTPS()))
}

func (g *G) Update(screen *ebiten.Image) error {
	g.mgr.Update(1.0/60.0, 800, 600)
	return nil
}

func (g *G) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}
