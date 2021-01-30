package renderer

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/inkyblackness/imgui-go/v2"
)

var keys = map[int]int{
	imgui.KeyTab:        int(ebiten.KeyTab),
	imgui.KeyLeftArrow:  int(ebiten.KeyLeft),
	imgui.KeyRightArrow: int(ebiten.KeyRight),
	imgui.KeyUpArrow:    int(ebiten.KeyUp),
	imgui.KeyDownArrow:  int(ebiten.KeyDown),
	imgui.KeyPageUp:     int(ebiten.KeyPageUp),
	imgui.KeyPageDown:   int(ebiten.KeyPageDown),
	imgui.KeyHome:       int(ebiten.KeyHome),
	imgui.KeyEnd:        int(ebiten.KeyEnd),
	imgui.KeyInsert:     int(ebiten.KeyInsert),
	imgui.KeyDelete:     int(ebiten.KeyDelete),
	imgui.KeyBackspace:  int(ebiten.KeyBackspace),
	imgui.KeySpace:      int(ebiten.KeySpace),
	imgui.KeyEnter:      int(ebiten.KeyEnter),
	imgui.KeyEscape:     int(ebiten.KeyEscape),
	imgui.KeyA:          int(ebiten.KeyA),
	imgui.KeyC:          int(ebiten.KeyC),
	imgui.KeyV:          int(ebiten.KeyV),
	imgui.KeyX:          int(ebiten.KeyX),
	imgui.KeyY:          int(ebiten.KeyY),
	imgui.KeyZ:          int(ebiten.KeyZ),
}

func sendInput(io *imgui.IO) {

	// Ebiten hides the LeftAlt RightAlt implementation (inside the uiDriver()), so
	// here only the left alt is sent
	if ebiten.IsKeyPressed(ebiten.KeyAlt) {
		io.KeyAlt(1, 0)
	} else {
		io.KeyAlt(0, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyShift) {
		io.KeyShift(1, 0)
	} else {
		io.KeyShift(0, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		io.KeyCtrl(1, 0)
	} else {
		io.KeyCtrl(0, 0)
	}
	// TODO: get KeySuper somehow (GLFW: KeyLeftSuper    = Key(343), R: 347)

	if chars := ebiten.InputChars(); len(chars) > 0 {
		io.AddInputCharacters(string(chars))
	}
	for _, iv := range keys {
		if inpututil.IsKeyJustPressed(ebiten.Key(iv)) {
			io.KeyPress(iv)
		}
		if inpututil.IsKeyJustReleased(ebiten.Key(iv)) {
			io.KeyRelease(iv)
		}
	}
}
