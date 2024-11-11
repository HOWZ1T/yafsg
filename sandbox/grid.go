package sandbox

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"yafsg/sandbox/particles"
)

type Grid struct {
	Width   int
	Height  int
	Cells   []int
	Texture rl.Texture2D
	Stats   particles.Stats
}

func NewGrid(width, height int) Grid {
	// create blank texture
	blankImg := rl.GenImageColor(width, height, rl.Black)
	tex := rl.LoadTextureFromImage(blankImg)
	rl.UnloadImage(blankImg)

	return Grid{Width: width, Height: height, Cells: make([]int, width*height), Texture: tex, Stats: particles.Stats{Empty: width * height}}
}

func (g *Grid) Index(x, y int) int {
	return y*g.Width + x
}

func (g *Grid) Get(x, y int) int {
	return g.Cells[g.Index(x, y)]
}

func (g *Grid) Set(x, y, val int) {
	g.Cells[g.Index(x, y)] = val
}

// Render draws the grid to the grid's texture.
func (g *Grid) Render() {
	pixelCols := make([]rl.Color, g.Width*g.Height)
	for i := 0; i < len(g.Cells); i++ {
		pixelCols[i] = particles.Color(g.Cells[i])
	}

	rl.UpdateTexture(g.Texture, pixelCols)
}

func (g *Grid) inBounds(x, y int) bool {
	return x >= 0 && x < g.Width && y >= 0 && y < g.Height
}

func (g *Grid) solidRules(cellParticle, x, y int) {
	xOffsets := []int{0, -1, 1} // checks in order, bottom, left, right
	yOffset := 1

	if !g.inBounds(x, y+yOffset) {
		return
	}

	for _, xOffset := range xOffsets {
		// skip out of bounds
		if !g.inBounds(x+xOffset, y+yOffset) {
			continue
		}

		// if cell below is empty, move down
		if g.Get(x+xOffset, y+1) == particles.Empty {
			g.Set(x+xOffset, y+1, cellParticle)
			g.Set(x, y, 0)
			return
		}
	}
}

func (g *Grid) liquidRules(cellParticle, x, y int) {
	xOffsets := []int{0, -1, 1} // checks in order: bottom, left, right
	yOffsets := []int{1, 0}     // checks in order: bottom row, same row

	for _, yOffset := range yOffsets {
		if yOffset == 0 {
			// side to side checks state
			// add randomness by randomizing the order of the side to side checks
			if rl.GetRandomValue(0, 1) == 0 {
				xOffsets = []int{0, 1, -1} // checks in order: bottom, right, left
			}
		}

		for _, xOffset := range xOffsets {
			targetX := x + xOffset
			targetY := y + yOffset
			if !g.inBounds(targetX, targetY) {
				continue // skip out of bounds
			}

			if targetX == x && targetY == y {
				continue // skip self
			}

			// if target cell is empty, move there
			if g.Get(targetX, targetY) == particles.Empty {
				g.Set(targetX, targetY, cellParticle)
				g.Set(x, y, 0)
				return
			}
		}
	}
}

// Tick simulates one step
func (g *Grid) Tick() error {
	for y := g.Height - 1; y >= 0; y-- {
		for x := 0; x < g.Width; x++ {
			cellParticle := g.Get(x, y)

			switch cellParticle {
			case particles.Empty:
				continue

			case particles.Sand:
				g.solidRules(cellParticle, x, y)

			case particles.Water:
				g.liquidRules(cellParticle, x, y)

			default:
				return fmt.Errorf("unknown particle type: %d", cellParticle)
			}
		}
	}

	return nil
}

func (g *Grid) UpdateStats() {
	g.Stats = particles.Stats{}
	for _, cell := range g.Cells {
		switch cell {
		case particles.Empty:
			g.Stats.Empty++

		case particles.Sand:
			g.Stats.Sand++

		case particles.Water:
			g.Stats.Water++

		default:
			g.Stats.Unknown++
		}
	}
}

func (g *Grid) Close() {
	rl.UnloadTexture(g.Texture)
}
