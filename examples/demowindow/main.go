//go:build example
// +build example

package main

import (
	"fmt"

	imgui "github.com/gabstv/cimgui-go"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Example with the main Demo window and ClipMask

func main() {
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	gg := &G{
		dscale:         ebiten.DeviceScaleFactor(),
		showDemoWindow: true,
	}

	ebiten.RunGame(gg)
}

type G struct {
	showDemoWindow bool
	dscale         float64
	retina         bool
	w, h           int
}

func (g *G) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.2f\nFPS: %.2f\n[C]lipMask: %t", ebiten.ActualTPS(), ebiten.ActualFPS(), ebimgui.ClipMask()), 10, 2)
	ebimgui.Draw(screen)
}

func (g *G) Update() error {
	ebimgui.Update(1.0 / 60.0)

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		ebimgui.SetClipMask(!ebimgui.ClipMask())
	}

	ebimgui.BeginFrame()
	defer ebimgui.EndFrame()

	imgui.Checkbox("Retina", &g.retina) // Edit bools storing our window open/close state

	imgui.Checkbox("Demo Window", &g.showDemoWindow) // Edit bools storing our window open/close state

	if g.showDemoWindow {
		imgui.ShowDemoWindow()
	}
	return nil
}

func lerp(a, b, t float64) float64 {
	return a*(1-t) + b*t
}

func (g *G) Layout(outsideWidth, outsideHeight int) (int, int) {
	if g.retina {
		m := ebiten.DeviceScaleFactor()
		g.w = int(float64(outsideWidth) * m)
		g.h = int(float64(outsideHeight) * m)
	} else {
		g.w = outsideWidth
		g.h = outsideHeight
	}
	ebimgui.SetDisplaySize(float32(g.w), float32(g.h))
	return g.w, g.h
}
