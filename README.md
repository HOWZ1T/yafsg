# YAFSG
Yet Another Falling Sand Game.

---

This project is a work in progress and is just an excuse to mess around with
implementing a falling sand simulation.

If it develops into anything more, I'll be sure to update this README.

### TODO
- [ ] add gravity rules for particle acceleration when falling
- [ ] add particle interaction rules, i.e. solids falls through liquids
- [ ] add fixed solids
- [ ] add dirty rectangles to control which areas need updates
- [ ] add multi-threading

---

## Controls
- `Left Mouse` - Place a particle.
- `Right Mouse [Hold]` - Pan the camera.
- `Mouse Wheel` - Zoom in/out.
- `Space` - Step the simulation if paused.
- `P` - Pause the simulation.
- `D` - Toggle debug mode.

---

## Building
This project uses raylib-go as the rendering backend.
Please ensure you have a _*64-bit c*_ compiler set for use in your go build toolchain (e.g.: `export CC=<path-to-64bit-compiler>`).

