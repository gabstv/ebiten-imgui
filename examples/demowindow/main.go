// +build example

package main

import (
	"fmt"

	"github.com/gabstv/ebiten-imgui/renderer"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/inkyblackness/imgui-go/v2"
)

// Example with the main Demo window and ClipMask

func main() {
	mgr := renderer.New(nil)

	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowResizable(true)

	gg := &G{
		mgr:    mgr,
		dscale: ebiten.DeviceScaleFactor(),
	}

	ebiten.RunGame(gg)
}

type G struct {
	mgr *renderer.Manager
	// demo members:
	showDemoWindow bool
	dscale         float64
	retina         bool
	w, h           int
}

func (g *G) Draw(screen *ebiten.Image) {
	g.mgr.BeginFrame()

	{
		imgui.Checkbox("Retina", &g.retina) // Edit bools storing our window open/close state

		imgui.Checkbox("Demo Window", &g.showDemoWindow) // Edit bools storing our window open/close state

		if g.showDemoWindow {
			imgui.ShowDemoWindow(&g.showDemoWindow)
		}
	}

	g.mgr.EndFrame(screen)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.2f", ebiten.CurrentTPS()), 10, 2)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[C]lipMask: %t", g.mgr.ClipMask), 10, 20)
}

func (g *G) Update() error {
	g.mgr.Update(1.0/60.0, float32(g.w), float32(g.h))
	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.mgr.ClipMask = !g.mgr.ClipMask
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
	return g.w, g.h
}
