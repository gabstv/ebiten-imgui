package renderer

import (
	"runtime"
	"sync/atomic"
	"unsafe"

	"github.com/hajimehoshi/ebiten"
	"github.com/inkyblackness/imgui-go/v2"
)

var nextTextureID int32

type GetCursorFn func() (x, y float32)

type Manager struct {
	Cache        map[imgui.TextureID]*ebiten.Image
	ctx          *imgui.Context
	cliptxt      string
	rawatlas     []uint8
	GetCursor    GetCursorFn
	SyncInputsFn func()
	SyncCursor   bool
	SyncInputs   bool
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
	Render(screen, imgui.RenderedDrawData(), m.Cache)
}

func New(fontAtlas *imgui.FontAtlas) *Manager {
	imctx := imgui.CreateContext(fontAtlas)
	m := &Manager{
		Cache:      make(map[imgui.TextureID]*ebiten.Image),
		ctx:        imctx,
		SyncCursor: true,
		SyncInputs: true,
	}
	runtime.SetFinalizer(m, (*Manager).onfinalize)
	// Build texture atlas
	io := imgui.CurrentIO()
	image := io.Fonts().TextureDataAlpha8()
	m.rawatlas = make([]uint8, image.Width*image.Height)
	for i := range m.rawatlas {
		m.rawatlas[i] = 255
	}
	image.Pixels = unsafe.Pointer(&m.rawatlas[0])
	id := atomic.AddInt32(&nextTextureID, 1)
	io.Fonts().SetTextureID(imgui.TextureID(id))

	m.setKeyMapping()

	return m
}

func NewWithContext(ctx *imgui.Context) *Manager {
	m := &Manager{
		Cache:      make(map[imgui.TextureID]*ebiten.Image),
		ctx:        ctx,
		SyncCursor: true,
		SyncInputs: true,
	}
	m.setKeyMapping()
	return m
}

type Renderer struct {
	Target *ebiten.Image
	Cache  map[imgui.TextureID]*ebiten.Image
}

func (r *Renderer) PreRender(clearColor [3]float32) {
	_ = r.Target.Clear()
}

func (r *Renderer) Render(drawData imgui.DrawData) {
	if r.Cache == nil {
		r.Cache = make(map[imgui.TextureID]*ebiten.Image)
	}
	Render(r.Target, drawData, r.Cache)
}
