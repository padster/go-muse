// Renders various data from a Muse headset to screen.
package ui

import (
	"log"
	"math"
	"runtime"
	"time"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

type Screen struct {
	width int
	height int
}

func NewScreen(width int, height int) *Screen {
	s := Screen {
		width,
		height,
	}
	return &s
}

func (s *Screen) Start() {
	runtime.LockOSThread()

	glfw.SetErrorCallback(func(err glfw.ErrorCode, desc string) {
		log.Fatalf("%v: %s\n", err, desc)
	})
	if !glfw.Init() {
		log.Fatalf("Can't init glfw!")
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Decorated, 1)
	glfw.WindowHint(glfw.Resizable, 0)

	window, err := glfw.CreateWindow(s.width, s.height, "Muse", nil, nil)
	if err != nil {
		log.Fatalf("CreateWindow failed: %s", err)
	}
	aw, ah := window.GetSize()
	if aw != s.width || ah != s.height {
		log.Fatalf("Window doesn't have the requested size: want %d,%d got %d,%d", 
			s.width, s.height, aw, ah)
	}
	window.MakeContextCurrent()

	start := time.Now()
	for !window.ShouldClose() {
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			break
		}
		gl.ClearColor(0, 0, 0, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.MatrixMode(gl.MODELVIEW_MATRIX)
		gl.LoadIdentity()
		gl.Translated(-1, -1, 0)
		gl.Scaled(2/float64(s.width), 2/float64(s.height), 1.0)
		gl.Begin(gl.LINE_STRIP)
		for i := 0; i < s.width; i++ {
			gl.Vertex2i(i, s.f(start, i))
		}
		gl.End()
		gl.Begin(gl.LINES)
		for i := 0; i < s.width; i += 8 {
			gl.Vertex2i(i, s.height/2)
			gl.Vertex2i(i, s.f(start, i))
		}
		gl.End()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

// TODO: Connect to channel rather than generate values here.
func (s *Screen) f(start time.Time, i int) int {
	theta := time.Since(start).Seconds() * 3.9
	x0 := float64(i) * 3 * math.Pi / float64(s.width)
	x0 += 0.4 * math.Sin(x0*2) * math.Cos(theta)
	x := x0 + 5*time.Since(start).Seconds()
	y := math.Cos(x0) + 0.5*math.Cos(x*2.432) + 0.25*math.Cos(x*4.123) + 0.125*math.Cos(x*8.3847)
	return int(y*float64(s.height)*0.3) + s.height/2
}