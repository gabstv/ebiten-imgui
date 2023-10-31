package renderer

import (
	"github.com/AllenDang/cimgui-go"
	"github.com/hajimehoshi/ebiten/v2"
)

type TextureCache interface {
	FontAtlasTextureID() cimgui.ImTextureID
	SetFontAtlasTextureID(id cimgui.ImTextureID)
	GetTexture(id cimgui.ImTextureID) *ebiten.Image
	SetTexture(id cimgui.ImTextureID, img *ebiten.Image)
	RemoveTexture(id cimgui.ImTextureID)
	ResetFontAtlasCache(filter ebiten.Filter)
}

type textureCache struct {
	fontAtlasID    cimgui.ImTextureID
	fontAtlasImage *ebiten.Image
	cache          map[cimgui.ImTextureID]*ebiten.Image
	dfilter        ebiten.Filter
}

var _ TextureCache = (*textureCache)(nil)

func (c *textureCache) getFontAtlas() *ebiten.Image {
	if c.fontAtlasImage == nil {
		pixels, width, height, outBytesPerPixel := cimgui.GetIO().GetFonts().GetTextureDataAsRGBA32() // call this to force imgui to build the font atlas cache
		c.fontAtlasImage = getTexture(pixels, width, height, outBytesPerPixel)
	}
	return c.fontAtlasImage
}

func (c *textureCache) FontAtlasTextureID() cimgui.ImTextureID {
	return c.fontAtlasID
}

func (c *textureCache) SetFontAtlasTextureID(id cimgui.ImTextureID) {
	c.fontAtlasID = id
	// c.fontAtlasImage = nil
}

func (c *textureCache) GetTexture(id cimgui.ImTextureID) *ebiten.Image {
	if id != c.fontAtlasID {
		if im, ok := c.cache[id]; ok {
			return im
		}
	}
	return c.getFontAtlas()
}

func (c *textureCache) SetTexture(id cimgui.ImTextureID, img *ebiten.Image) {
	c.cache[id] = img
}

func (c *textureCache) RemoveTexture(id cimgui.ImTextureID) {
	delete(c.cache, id)
}

func (c *textureCache) ResetFontAtlasCache(filter ebiten.Filter) {
	c.fontAtlasImage = nil
	c.dfilter = filter
}

func NewCache() TextureCache {
	return &textureCache{
		fontAtlasID:    1,
		cache:          make(map[cimgui.ImTextureID]*ebiten.Image),
		fontAtlasImage: nil,
	}
}
