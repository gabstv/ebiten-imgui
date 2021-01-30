package renderer

import (
	"runtime"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/inkyblackness/imgui-go/v2"
)

//var nextTextureID int32

type GetCursorFn func() (x, y float32)

type Manager struct {
	Filter       ebiten.Filter
	Cache        TextureCache
	ctx          *imgui.Context
	cliptxt      string
	rawatlas     []uint8
	GetCursor    GetCursorFn
	SyncInputsFn func()
	SyncCursor   bool
	SyncInputs   bool
	lmask        *ebiten.Image
	ClipMask     bool
}

func (m *Manager) onfinalize() {
	runtime.SetFinalizer(m, nil)
	m.ctx.Destroy()
}

func (m *Manager) setKeyMapping() {
	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	io := imgui.CurrentIO()
	for imguiKey, nativeKey := range keys {
		io.KeyMap(imguiKey, nativeKey)
	}
}

// Text implements imgui clipboard
func (m *Manager) Text() (string, error) {
	return m.cliptxt, nil
}

// SetText implements imgui clipboard
func (m *Manager) SetText(text string) {
	m.cliptxt = text
}

func (m *Manager) Update(delta, winWidth, winHeight float32) {
	io := imgui.CurrentIO()
	io.SetDisplaySize(imgui.Vec2{X: winWidth, Y: winHeight})
	io.SetDeltaTime(delta)
	if m.SyncCursor {
		if m.GetCursor != nil {
			x, y := m.GetCursor()
			io.SetMousePosition(imgui.Vec2{X: x, Y: y})
		} else {
			mx, my := ebiten.CursorPosition()
			io.SetMousePosition(imgui.Vec2{X: float32(mx), Y: float32(my)})
		}
		io.SetMouseButtonDown(0, ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft))
		io.SetMouseButtonDown(1, ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight))
		io.SetMouseButtonDown(2, ebiten.IsMouseButtonPressed(ebiten.MouseButtonMiddle))
		xoff, yoff := ebiten.Wheel()
		io.AddMouseWheelDelta(float32(xoff), float32(yoff))
	}
	if m.SyncInputs {
		if m.SyncInputsFn != nil {
			m.SyncInputsFn()
		} else {
			sendInput(&io)
		}
	}
}

func (m *Manager) BeginFrame() {
	imgui.NewFrame()
}

func (m *Manager) EndFrame(screen *ebiten.Image) {
	imgui.Render()
	if m.ClipMask {
		if m.lmask == nil {
			w, h := screen.Size()
			m.lmask = ebiten.NewImage(w, h)
		} else {
			w1, h1 := screen.Size()
			w2, h2 := m.lmask.Size()
			if w1 != w2 || h1 != h2 {
				m.lmask.Dispose()
				m.lmask = ebiten.NewImage(w1, h1)
			}
		}
		RenderMasked(screen, m.lmask, imgui.RenderedDrawData(), m.Cache, m.Filter)
	} else {
		Render(screen, imgui.RenderedDrawData(), m.Cache, m.Filter)
	}
}

func New(fontAtlas *imgui.FontAtlas) *Manager {
	imctx := imgui.CreateContext(fontAtlas)
	m := &Manager{
		Cache:      NewCache(),
		ctx:        imctx,
		SyncCursor: true,
		SyncInputs: true,
		ClipMask:   true,
	}
	runtime.SetFinalizer(m, (*Manager).onfinalize)
	// Build texture atlas
	io := imgui.CurrentIO()
	_ = io.Fonts().TextureDataRGBA32() // call this to force imgui to build the font atlas cache
	io.Fonts().SetTextureID(1)
	m.Cache.SetFontAtlasTextureID(1)

	m.setKeyMapping()

	return m
}

func NewWithContext(ctx *imgui.Context) *Manager {
	m := &Manager{
		Cache:      NewCache(),
		ctx:        ctx,
		SyncCursor: true,
		SyncInputs: true,
		ClipMask:   true,
	}
	m.setKeyMapping()
	return m
}
