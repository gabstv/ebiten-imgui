// Package imcolor is a utility package to convert between color.Color and imgui.Vec4.
package imcolor

import (
	"image/color"

	"github.com/AllenDang/cimgui-go"
)

func ToVec4(c color.Color) cimgui.ImVec4 {
	r32, g32, b32, a32 := c.RGBA()
	return cimgui.ImVec4{
		X: float32(r32) / 0xffff,
		Y: float32(g32) / 0xffff,
		Z: float32(b32) / 0xffff,
		W: float32(a32) / 0xffff,
	}
}

func ToColor(c cimgui.ImVec4) color.RGBA64 {
	return color.RGBA64{
		R: uint16(c.X * 0xffff),
		G: uint16(c.Y * 0xffff),
		B: uint16(c.Z * 0xffff),
		A: uint16(c.W * 0xffff),
	}
}
