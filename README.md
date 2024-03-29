Dear ImGui for Ebitengine
=========================

A renderer of [cimgui-go](https://github.com/gabstv/cimgui-go) for [Ebitengine](https://github.com/hajimehoshi/ebiten/)!

![preview](assets/preview.png)

Usage:

```go
package main

import (
	"fmt"
	"image/color"

	"github.com/gabstv/ebiten-imgui/v3/imcolor"
	ebimgui "github.com/gabstv/ebiten-imgui/v3"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	imgui "github.com/gabstv/cimgui-go"
)

func main() {
	gg := &G{
		name:       "Hello, Dear ImGui",
		clearColor: [3]float32{0,0,0},
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle(gg.name)

	ebiten.RunGame(gg)
}

type G struct {
	clearColor [3]float32
	floatVal   float32
	counter    int
	name       string
}

func (g *G) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{uint8(g.clearColor[0] * 255), uint8(g.clearColor[1] * 255), uint8(g.clearColor[2] * 255), 255})
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %.2f", ebiten.CurrentTPS()))
	ebimgui.Draw(screen)
}

func InputText(label string, buf *string) bool {
	return imgui.InputTextWithHint(label, "", buf, 0, nil)
}

func (g *G) Update() error {
	ebimgui.Update(1.0 / 60.0)
	ebimgui.BeginFrame()
	defer ebimgui.EndFrame()

	imgui.Text("ภาษาไทย测试조선말")                        // To display these, you'll need to register a compatible font
	imgui.Text("Hello, world!")                       // Display some text
	imgui.SliderFloat("float", &g.floatVal, 0.0, 1.0) // Edit 1 float using a slider from 0.0f to 1.0f
	imgui.ColorEdit3("clear color", &g.clearColor)     // Edit 3 floats representing a color

	//imgui.Checkbox("Demo Window", &showDemoWindow) // Edit bools storing our window open/close state
	//imgui.Checkbox("Go Demo Window", &showGoDemoWindow)
	//imgui.Checkbox("Another Window", &showAnotherWindow)

	if imgui.Button("Button") { // Buttons return true when clicked (most widgets return true when edited/activated)
		g.counter++
	}
	imgui.SameLine()
	imgui.Text(fmt.Sprintf("counter = %d", g.counter))

	if InputText("Window title", &g.name) {
		ebiten.SetWindowTitle(g.name)
	}

	xcol := imcolor.ToVec4(color.RGBA{
		R: 0xFF,
		G: 0x00,
		B: 0xFF,
		A: 0x99,
	})

	imgui.PushStyleColor_Vec4(imgui.ImGuiCol_Text, xcol)
	imgui.Text(fmt.Sprintf("fps = %f", ebiten.CurrentFPS()))
	imgui.PopStyleColor()
	return nil
}

func (g *G) Layout(outsideWidth, outsideHeight int) (int, int) {
	ebimgui.SetDisplaySize(float32(800), float32(600))
	return 800, 600
}
```
