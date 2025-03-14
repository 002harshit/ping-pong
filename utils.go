package main

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawRectCenter(cx, cy, w, h int, col color.RGBA) {
	rl.DrawRectangle(int32(cx-w/2), int32(cy-h/2), int32(w), int32(h), col)
}
func frame_time() float32 {
	return rl.GetFrameTime() * 60
}
func isCollidingCircleRec(b Ball, radius float32, p Entity, size rl.Vector2) bool {

	return rl.CheckCollisionCircleRec(b.pos, radius, rl.Rectangle{X: p.pos.X, Y: p.pos.Y, Width: size.X, Height: size.Y})
}
func DrawTextCenter(message string, cx, cy, fontsize int, col color.RGBA) {
	size := rl.MeasureTextEx(rl.GetFontDefault(), message, float32(fontsize), float32(fontsize/10))
	rl.DrawText(message, int32(float32(cx)-size.X/2), int32(float32(cy)-size.Y/2), int32(fontsize), rl.Black)
}
