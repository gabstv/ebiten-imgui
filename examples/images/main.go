//go:build example
// +build example

package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	imgui "github.com/gabstv/cimgui-go"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var exampleImage *ebiten.Image

var (
	myImageIDRef int = 10
)

func main() {
	ebiten.SetWindowSize(800, 600)

	gg := &G{}

	exampleImage, _, err := ebitenutil.NewImageFromFile("example.png")
	if err != nil {
		log.Fatal(err)
	}
	ebimgui.GlobalManager().Cache.SetTexture(imgui.TextureID(&myImageIDRef), exampleImage) // Texture ID 10 will contain this example image

	ebiten.RunGame(gg)
}

type G struct{}

func (g *G) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.2f", ebiten.CurrentTPS()))
	ebimgui.Draw(screen)
}

// cimgui-go doesn't have imgui.Image with sensible default yet
// you can implement it like this
func Image(tid imgui.TextureID, size imgui.Vec2) {
	uv0 := imgui.NewVec2(0, 0)
	uv1 := imgui.NewVec2(1, 1)
	border_col := imgui.NewVec4(0, 0, 0, 0)
	tint_col := imgui.NewVec4(1, 1, 1, 1)

	imgui.ImageV(tid, size, uv0, uv1, tint_col, border_col)
}

func (g *G) Update() error {
	ebimgui.Update(1.0 / 60.0)
	ebimgui.BeginFrame()
	imgui.Text("Hello, images!")
	Image(imgui.TextureID(&myImageIDRef), imgui.Vec2{X: 64, Y: 64})
	ebimgui.EndFrame()
	return nil
}

func (g *G) Layout(outsideWidth, outsideHeight int) (int, int) {
	ebimgui.SetDisplaySize(float32(800), float32(600))
	return 800, 600
}
