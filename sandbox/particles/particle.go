package particles

import rl "github.com/gen2brain/raylib-go/raylib"

const (
	Empty = iota
	Concrete
	Sand
	Water
)

func Color(particleType int) rl.Color {
	switch particleType {
	case Empty:
		return rl.DarkPurple

	case Concrete:
		return rl.Gray

	case Sand:
		return rl.Beige

	case Water:
		return rl.Blue

	default:
		return rl.Magenta
	}
}

func Next(i int) int {
	return (i + 1) % 4
}

func Prev(i int) int {
	return (i + 3) % 4
}

func Name(i int) string {
	switch i {
	case Empty:
		return "Empty"

	case Concrete:
		return "Concrete"

	case Sand:
		return "Sand"

	case Water:
		return "Water"

	default:
		return "Unknown"
	}
}
