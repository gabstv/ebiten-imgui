package renderer

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/inkyblackness/imgui-go/v2"
)

/*
// Renderer covers rendering imgui draw data.
type Renderer interface {
	// PreRender causes the display buffer to be prepared for new output.
	PreRender(clearColor [3]float32)
	// Render draws the provided imgui draw data.
	Render(displaySize [2]float32, framebufferSize [2]float32, drawData imgui.DrawData)
}
*/

type Renderer struct {
	Target *ebiten.Image
	Cache  map[imgui.TextureID]*ebiten.Image
}

func (r *Renderer) PreRender(clearColor [3]float32) {
	_ = r.Target.Clear()
}

func (r *Renderer) Render(displaySize [2]float32, framebufferSize [2]float32, drawData imgui.DrawData) {
	if r.Cache == nil {
		r.Cache = make(map[imgui.TextureID]*ebiten.Image)
	}
	Render(r.Target, displaySize, framebufferSize, drawData, r.Cache)
}
