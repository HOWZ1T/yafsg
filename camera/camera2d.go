package camera

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"yafsg/mathutils"
)

type Camera2D struct {
	scale         float32
	zoom          Zoom
	boundedOffset mathutils.BoundedOffset2
}

func NewCamera2D(scale float32, zoom Zoom, boundedOffset mathutils.BoundedOffset2) Camera2D {
	return Camera2D{scale, zoom, boundedOffset}
}

func (cam2d *Camera2D) Zoom(dir float32, scrWidth, scrHeight, playAreaWidth, playAreaHeight int) {
	scrW := float32(scrWidth)
	scrH := float32(scrHeight)
	playAreaW := float32(playAreaWidth)
	playAreaH := float32(playAreaHeight)

	cam2d.scale = rl.Clamp(cam2d.scale+cam2d.zoom.Speed*dir, cam2d.zoom.Min, cam2d.zoom.Max)
	cam2d.boundedOffset.Min = rl.NewVector2(scrW-playAreaW*cam2d.scale, scrH-playAreaH*cam2d.scale)
	cam2d.boundedOffset.Max = rl.NewVector2(0, 0)
	cam2d.boundedOffset.Offset = rl.Vector2Clamp(cam2d.boundedOffset.Offset, cam2d.boundedOffset.Min, cam2d.boundedOffset.Max)
}

func (cam2d *Camera2D) ZoomAbout(scrX, scrY int) {
	// TODO
	panic("Not implemented")
}

func (cam2d *Camera2D) Offset() rl.Vector2 {
	return cam2d.boundedOffset.Offset
}

func (cam2d *Camera2D) Scale() float32 {
	return cam2d.scale
}

func (cam2d *Camera2D) Pan(dir rl.Vector2) {
	cam2d.boundedOffset.Offset = rl.Vector2Add(cam2d.boundedOffset.Offset, dir)
	cam2d.boundedOffset.Offset = rl.Vector2Clamp(cam2d.boundedOffset.Offset, cam2d.boundedOffset.Min, cam2d.boundedOffset.Max)
}

func (cam2d *Camera2D) ScrToWorldPixel(scrPos rl.Vector2, playAreaWidth, playAreaHeight int) rl.Vector2 {
	scrPos = rl.Vector2Subtract(scrPos, cam2d.Offset())
	scrPos = rl.Vector2Divide(scrPos, rl.NewVector2(cam2d.Scale(), cam2d.Scale()))
	scrPos = rl.Vector2Clamp(scrPos, rl.NewVector2(0, 0), rl.NewVector2(float32(playAreaWidth-1), float32(playAreaHeight-1)))
	worldPixelX := math.Floor(float64(scrPos.X))
	worldPixelY := math.Floor(float64(scrPos.Y))

	return rl.NewVector2(float32(worldPixelX), float32(worldPixelY))
}

func (cam2d *Camera2D) WorldPixelToScr(worldPixel rl.Vector2) rl.Vector2 {
	scrPixel := rl.NewVector2(worldPixel.X, worldPixel.Y)
	scrPixel = rl.Vector2Multiply(scrPixel, rl.NewVector2(cam2d.Scale(), cam2d.Scale()))
	scrPixel = rl.Vector2Add(scrPixel, cam2d.Offset())

	return scrPixel
}
