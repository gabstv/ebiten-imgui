package renderer

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/gabstv/ebiten-imgui/internal/native"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/inkyblackness/imgui-go/v2"
)

var pixelimg *ebiten.Image

// struct ImDrawVert
// {
//     ImVec2  pos; // 2 floats
//     ImVec2  uv; // 2 floats
//     ImU32   col; // uint32
// };

type cVec2x32 struct {
	X float32
	Y float32
}

type cImDrawVertx32 struct {
	Pos cVec2x32
	UV  cVec2x32
	Col uint32
}

type cVec2x64 struct {
	X float64
	Y float64
}

type cImDrawVertx64 struct {
	Pos cVec2x64
	UV  cVec2x64
	Col uint32
}

func getVertices(vbuf unsafe.Pointer, vblen, vsize, offpos, offuv, offcol int) []ebiten.Vertex {
	if native.SzFloat() == 4 {
		return getVerticesx32(vbuf, vblen, vsize, offpos, offuv, offcol)
	}
	if native.SzFloat() == 8 {
		return getVerticesx64(vbuf, vblen, vsize, offpos, offuv, offcol)
	}
	panic("invalid char size")
}

func getVerticesx32(vbuf unsafe.Pointer, vblen, vsize, offpos, offuv, offcol int) []ebiten.Vertex {
	n := vblen / vsize
	vertices := make([]ebiten.Vertex, 0, vblen/vsize)
	rawverts := (*[1 << 28]cImDrawVertx32)(vbuf)[:n:n]
	for i := 0; i < n; i++ {
		c0 := rawverts[i].Col
		c00 := uint8(c0 & 0xFF)
		c01 := (c0 >> 8) & 0xFF
		c02 := (c0 >> 16) & 0xFF
		c03 := (c0 >> 24) & 0xFF
		_, _, _, _ = c00, c01, c02, c03
		vertices = append(vertices, ebiten.Vertex{
			SrcX:   rawverts[i].UV.X,
			SrcY:   rawverts[i].UV.Y,
			DstX:   rawverts[i].Pos.X,
			DstY:   rawverts[i].Pos.Y,
			ColorR: float32(rawverts[i].Col&0xFF) / 255,
			ColorG: float32(rawverts[i].Col>>8&0xFF) / 255,
			ColorB: float32(rawverts[i].Col>>16&0xFF) / 255,
			ColorA: float32(rawverts[i].Col>>24&0xFF) / 255,
		})
	}
	return vertices
}

func getVerticesx64(vbuf unsafe.Pointer, vblen, vsize, offpos, offuv, offcol int) []ebiten.Vertex {
	n := vblen / vsize
	vertices := make([]ebiten.Vertex, 0, vblen/vsize)
	rawverts := (*[1 << 28]cImDrawVertx64)(vbuf)[:n:n]
	for i := 0; i < n; i++ {
		vertices = append(vertices, ebiten.Vertex{
			SrcX:   float32(rawverts[i].UV.X),
			SrcY:   float32(rawverts[i].UV.Y),
			DstX:   float32(rawverts[i].Pos.X),
			DstY:   float32(rawverts[i].Pos.Y),
			ColorR: float32(rawverts[i].Col&0xFF) / 255,
			ColorG: float32(rawverts[i].Col>>8&0xFF) / 255,
			ColorB: float32(rawverts[i].Col>>16&0xFF) / 255,
			ColorA: float32(rawverts[i].Col>>24&0xFF) / 255,
		})
	}
	return vertices
}

func lerp(a, b int, t float32) float32 {
	return float32(a)*(1-t) + float32(b)*t
}

func vcopy(v []ebiten.Vertex) []ebiten.Vertex {
	cl := make([]ebiten.Vertex, len(v))
	copy(cl, v)
	return cl
}

func vmultiply(v, vbuf []ebiten.Vertex, bmin, bmax image.Point) {
	for i := range vbuf {
		vbuf[i].SrcX = lerp(bmin.X, bmax.X, v[i].SrcX)
		vbuf[i].SrcY = lerp(bmin.Y, bmax.Y, v[i].SrcY)
	}
}

func getTexture(tex *imgui.RGBA32Image) *ebiten.Image {
	n := tex.Width * tex.Height
	pix := (*[1 << 28]uint8)(tex.Pixels)[: n*4 : n*4]
	img := ebiten.NewImage(tex.Width, tex.Height)
	img.ReplacePixels(pix)
	return img
}

func getIndices(ibuf unsafe.Pointer, iblen, isize int) []uint16 {
	n := iblen / isize
	switch isize {
	case 2:
		// direct conversion (without a data copy)
		//TODO: document the size limit (?) this fits 268435456 bytes
		// https://github.com/golang/go/wiki/cgo#turning-c-arrays-into-go-slices
		return (*[1 << 28]uint16)(ibuf)[:n:n]
	case 4:
		slc := make([]uint16, n)
		for i := 0; i < n; i++ {
			slc[i] = uint16(*(*uint32)(unsafe.Pointer(uintptr(ibuf) + uintptr(i*isize))))
		}
		return slc
	case 8:
		slc := make([]uint16, n)
		for i := 0; i < n; i++ {
			slc[i] = uint16(*(*uint64)(unsafe.Pointer(uintptr(ibuf) + uintptr(i*isize))))
		}
		return slc
	default:
		panic(fmt.Sprint("byte size", isize, "not supported"))
	}
	return nil
}

// Render the ImGui drawData into the target *ebiten.Image
func Render(target *ebiten.Image, drawData imgui.DrawData, txcache TextureCache, dfilter ebiten.Filter) {
	render(target, nil, drawData, txcache, dfilter)
}

// RenderMasked renders the ImGui drawData into the target *ebiten.Image with ebiten.CompositeModeCopy for masking
func RenderMasked(target *ebiten.Image, mask *ebiten.Image, drawData imgui.DrawData, txcache TextureCache, dfilter ebiten.Filter) {
	render(target, mask, drawData, txcache, dfilter)
}

func render(target *ebiten.Image, mask *ebiten.Image, drawData imgui.DrawData, txcache TextureCache, dfilter ebiten.Filter) {
	targetw, targeth := target.Size()
	if !drawData.Valid() {
		return
	}

	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.VertexBufferLayout()
	indexSize := imgui.IndexBufferLayout()

	opt := &ebiten.DrawTrianglesOptions{
		Filter: dfilter,
	}
	var opt2 *ebiten.DrawImageOptions
	if mask != nil {
		opt2 = &ebiten.DrawImageOptions{
			CompositeMode: ebiten.CompositeModeSourceOver,
		}
	}

	for _, clist := range drawData.CommandLists() {
		var indexBufferOffset int
		vertexBuffer, vertexLen := clist.VertexBuffer()
		indexBuffer, indexLen := clist.IndexBuffer()
		vertices := getVertices(vertexBuffer, vertexLen, vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol)
		vbuf := vcopy(vertices)
		indices := getIndices(indexBuffer, indexLen, indexSize)
		for _, cmd := range clist.Commands() {
			ecount := cmd.ElementCount()
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(clist)
			} else {
				clipRect := cmd.ClipRect()
				texid := cmd.TextureID()
				tx := txcache.GetTexture(texid)
				// if _, ok := txcache[texid]; !ok {
				// 	if texid == 1 {
				// 		txcache[texid] = getTexture(imgui.CurrentIO().Fonts().TextureDataRGBA32(), dfilter)
				// 	}
				// }
				// tx := txcache[texid]
				vmultiply(vertices, vbuf, tx.Bounds().Min, tx.Bounds().Max)
				if mask == nil || (clipRect.X == 0 && clipRect.Y == 0 && clipRect.Z == float32(targetw) && clipRect.W == float32(targeth)) {
					target.DrawTriangles(vbuf, indices[indexBufferOffset:indexBufferOffset+ecount], tx, opt)
				} else {
					mask.Clear()
					opt2.GeoM.Reset()
					opt2.GeoM.Translate(float64(clipRect.X), float64(clipRect.Y))
					mask.DrawTriangles(vbuf, indices[indexBufferOffset:indexBufferOffset+ecount], tx, opt)
					target.DrawImage(mask.SubImage(image.Rectangle{
						Min: image.Pt(int(clipRect.X), int(clipRect.Y)),
						Max: image.Pt(int(clipRect.Z), int(clipRect.W)),
					}).(*ebiten.Image), opt2)
				}
			}
			indexBufferOffset += ecount
		}
	}
}
