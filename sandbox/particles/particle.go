package particles

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	Empty = 0
	Sand  = 1
	Water = 2
)

func Color(particleType int) rl.Color {
	switch particleType {
	case Empty:
		return rl.DarkPurple

	case Sand:
		return rl.Beige

	case Water:
		return rl.Blue

	default:
		return rl.Magenta
	}
}

func Next(i int) int {
	return (i + 1) % 3
}

func Prev(i int) int {
	return (i + 2) % 3
}

func Name(i int) string {
	switch i {
	case Empty:
		return "Empty"

	case Sand:
		return "Sand"

	case Water:
		return "Water"

	default:
		return "Unknown"
	}
}
