package mathutils

import rl "github.com/gen2brain/raylib-go/raylib"

type BoundedOffset2 struct {
	Offset rl.Vector2
	Min    rl.Vector2
	Max    rl.Vector2
}
