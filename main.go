package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
	"yafsg/camera"
	"yafsg/mathutils"
	"yafsg/sandbox"
	"yafsg/sandbox/particles"
)

// TODO decouple sim from frame rate
func main() {
	width := 1920
	height := 1080
	playAreaWidth := 1920 / 4
	playAreaHeight := 1080 / 4

	rl.InitWindow(int32(width), int32(height), "yafsg")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	grid := sandbox.NewGrid(playAreaWidth, playAreaHeight)
	defer grid.Close()

	cam2d := camera.NewCamera2D(
		4,
		camera.Zoom{Speed: 0.5, Min: 4, Max: 40},
		mathutils.BoundedOffset2{
			Offset: rl.NewVector2(0, 0),
			Min:    rl.NewVector2(0, 0),
			Max:    rl.NewVector2(0, 0),
		},
	)

	isSimPaused := true
	cursorParticle := particles.Sand
	cursorSize := 5
	isDebug := false
	for !rl.WindowShouldClose() {
		// update
		// find cell under mouse
		mousePos := rl.GetMousePosition()

		// find the world pixel
		mouseWorldPixel := cam2d.ScrToWorldPixel(mousePos, playAreaWidth, playAreaHeight)

		mouseWheel := rl.GetMouseWheelMove()
		if mouseWheel != 0 {
			cam2d.Zoom(mouseWheel, width, height, playAreaWidth, playAreaHeight)
		}

		if rl.IsMouseButtonDown(rl.MouseRightButton) {
			mouseDelta := rl.GetMouseDelta()
			cam2d.Pan(mouseDelta)
		}

		if rl.IsKeyPressed(rl.KeyP) {
			isSimPaused = !isSimPaused
		}

		if rl.IsKeyPressed(rl.KeySpace) && isSimPaused {
			err := grid.Tick()
			if err != nil {
				// print error and close
				fmt.Printf("sim tick error: %v", err)
				break
			}
		}

		if rl.IsKeyPressed(rl.KeyRight) {
			cursorParticle = particles.Next(cursorParticle)
		}

		if rl.IsKeyPressed(rl.KeyLeft) {
			cursorParticle = particles.Prev(cursorParticle)
		}

		if rl.IsKeyDown(rl.KeyUp) {
			cursorSize++
		}

		if rl.IsKeyDown(rl.KeyDown) {
			cursorSize--
			if cursorSize < 2 {
				cursorSize = 2
			}
		}

		if rl.IsKeyPressed(rl.KeyD) {
			isDebug = !isDebug
		}

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			// draw cells
			for y := -cursorSize / 2; y < cursorSize/2; y++ {
				for x := -cursorSize / 2; x < cursorSize/2; x++ {
					targetPixelX := int(mouseWorldPixel.X) + x
					targetPixelY := int(mouseWorldPixel.Y) + y

					if !grid.InBounds(targetPixelX, targetPixelY) {
						continue
					}
					grid.Set(targetPixelX, targetPixelY, cursorParticle)
				}
			}
		}

		processedTicks := 0
		if !isSimPaused {
			ticks := 2
			t0 := rl.GetTime()
			for i := 0; i < ticks; i++ {
				err := grid.Tick()
				if err != nil {
					// print error and close
					fmt.Printf("sim tick error: %v", err)
					break
				}
				processedTicks++

				// check if taking longer than 16 ms
				t1 := rl.GetTime()

				if t1-t0 > 0.016 {
					break
				}
			}
		}

		if isDebug {
			grid.UpdateStats()
		}

		// render
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		grid.Render()
		rl.DrawTextureEx(grid.Texture, cam2d.Offset(), 0, cam2d.Scale(), rl.White)

		// draw highlight around cell under mouse
		cursorColor := particles.Color(cursorParticle)
		cursorName := particles.Name(cursorParticle)
		if cursorParticle == particles.Empty {
			cursorColor = rl.Pink
			cursorName = "Eraser"
		}

		cursorHighlightTopLeft := cam2d.WorldPixelToScr(rl.Vector2Add(mouseWorldPixel, rl.Vector2{X: float32(math.Floor(float64(-cursorSize / 2))), Y: float32(math.Floor(float64(-cursorSize / 2)))}))
		curHighlightBotRight := cam2d.WorldPixelToScr(rl.Vector2Add(mouseWorldPixel, rl.Vector2{X: float32(math.Floor(float64(cursorSize / 2))), Y: float32(math.Floor(float64(cursorSize / 2)))}))
		cursorHighlightSize := rl.Vector2Subtract(curHighlightBotRight, cursorHighlightTopLeft)
		cursorHighlightRec := rl.Rectangle{
			X:      cursorHighlightTopLeft.X,
			Y:      cursorHighlightTopLeft.Y,
			Width:  cursorHighlightSize.X,
			Height: cursorHighlightSize.Y,
		}
		rl.DrawRectangleLinesEx(
			cursorHighlightRec,
			2,
			rl.Gold,
		)
		rl.DrawText(cursorName, int32(mousePos.X+25), int32(mousePos.Y-25), 20, cursorColor)

		if isSimPaused {
			// draw pause icon which is two verticle lines
			rl.DrawRectangle(int32(mousePos.X-25), int32(mousePos.Y-25), 5, 20, rl.Red)
			rl.DrawRectangle(int32(mousePos.X-35), int32(mousePos.Y-25), 5, 20, rl.Red)
		}

		// render debug stats
		if isDebug {
			rl.DrawText("Stats", 10, 10, 20, rl.White)
			rl.DrawText(fmt.Sprintf("Empty: %d", grid.Stats.Empty), 10, 40, 20, rl.White)
			rl.DrawText(fmt.Sprintf("Concrete: %d", grid.Stats.Concrete), 10, 70, 20, rl.White)
			rl.DrawText(fmt.Sprintf("Sand: %d", grid.Stats.Sand), 10, 100, 20, rl.White)
			rl.DrawText(fmt.Sprintf("Water: %d", grid.Stats.Water), 10, 130, 20, rl.White)
			rl.DrawText(fmt.Sprintf("Unknown: %d", grid.Stats.Unknown), 10, 160, 20, rl.White)
			rl.DrawText(fmt.Sprintf("Total: %d", grid.Stats.Total()), 10, 190, 20, rl.White)
			rl.DrawText(fmt.Sprintf("Processed Ticks: %d", processedTicks), 10, 220, 20, rl.White)
			rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), 10, 250, 20, rl.White)
		}

		rl.EndDrawing()
	}
}
