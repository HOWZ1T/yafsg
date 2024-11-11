package main

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"yafsg/camera"
	"yafsg/mathutils"
	"yafsg/sandbox"
	"yafsg/sandbox/particles"
)

// TODO decouple sim from frame rate
func main() {
	width := 1920
	height := 1080
	playAreaWidth := 192
	playAreaHeight := 108

	rl.InitWindow(int32(width), int32(height), "yafsg")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	grid := sandbox.NewGrid(playAreaWidth, playAreaHeight)
	defer grid.Close()

	cam2d := camera.NewCamera2D(
		10,
		camera.Zoom{Speed: 0.5, Min: 10, Max: 40},
		mathutils.BoundedOffset2{
			Offset: rl.NewVector2(0, 0),
			Min:    rl.NewVector2(0, 0),
			Max:    rl.NewVector2(0, 0),
		},
	)

	isSimPaused := true
	cursorParticle := particles.Sand
	isDebug := false
	for !rl.WindowShouldClose() {
		// update
		// find cell under mouse
		mousePos := rl.GetMousePosition()

		// find the world pixel
		mouseWorldPixel := cam2d.ScrToWorldPixel(mousePos, playAreaWidth, playAreaHeight)

		// find the top left of the world pixel on the screen
		mouseWorldPixelScr := cam2d.WorldPixelToScr(mouseWorldPixel)

		mouseWheel := rl.GetMouseWheelMove()
		if mouseWheel != 0 {
			cam2d.Zoom(mouseWheel, width, height, playAreaWidth, playAreaHeight)
		}

		if rl.IsMouseButtonDown(rl.MouseRightButton) {
			mouseDelta := rl.GetMouseDelta()
			cam2d.Pan(mouseDelta)
		}

		if rl.IsMouseButtonDown(rl.MouseLeftButton) {
			grid.Set(int(mouseWorldPixel.X), int(mouseWorldPixel.Y), cursorParticle)
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

		if rl.IsKeyPressed(rl.KeyD) {
			isDebug = !isDebug
		}

		if !isSimPaused {
			err := grid.Tick()
			if err != nil {
				// print error and close
				fmt.Printf("sim tick error: %v", err)
				break
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
		rl.DrawRectangleV(mouseWorldPixelScr, rl.NewVector2(cam2d.Scale(), cam2d.Scale()), cursorColor)
		rl.DrawText(cursorName, int32(mousePos.X+25), int32(mousePos.Y-25), 20, cursorColor)

		if isSimPaused {
			// draw pause icon which is two verticle lines
			rl.DrawRectangle(int32(mousePos.X-25), int32(mousePos.Y-25), 5, 20, rl.Red)
			rl.DrawRectangle(int32(mousePos.X-35), int32(mousePos.Y-25), 5, 20, rl.Red)
		}

		// render particle stats
		if isDebug {
			rl.DrawText("Stats", 10, 10, 20, rl.White)
			rl.DrawText(fmt.Sprintf("Empty: %d", grid.Stats.Empty), 10, 40, 20, rl.White)
			rl.DrawText(fmt.Sprintf("Sand: %d", grid.Stats.Sand), 10, 70, 20, rl.White)
			rl.DrawText(fmt.Sprintf("Water: %d", grid.Stats.Water), 10, 100, 20, rl.White)
			rl.DrawText(fmt.Sprintf("Unknown: %d", grid.Stats.Unknown), 10, 130, 20, rl.White)
			rl.DrawText(fmt.Sprintf("Total: %d", grid.Stats.Total()), 10, 160, 20, rl.White)
		}

		rl.EndDrawing()
	}
}