// Renders various data from a Muse headset to screen.
package ui

import (
	"log"
	"runtime"

	"github.com/padster/go-muse/util"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

type Screen struct {
	width int
	height int
	buffer *util.Buffer
}

// NewScreen creates a new output screen of a given size.
func NewScreen(width int, height int) *Screen {
	s := Screen {
		width,
		height,
		util.NewBuffer(width),
	}
	return &s
}

// Render starts rendering a channel of waves samples to screen.
func (s *Screen) Render(values <-chan float64, sampleRate int) {
	runtime.LockOSThread()

	glfw.SetErrorCallback(func(err glfw.ErrorCode, desc string) {
		log.Fatalf("%v: %s\n", err, desc)
	})
	if !glfw.Init() {
		log.Fatalf("Can't init glfw!")
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(s.width, s.height, "Muse", nil, nil)
	if err != nil {
		log.Fatalf("CreateWindow failed: %s", err)
	}
	if aw, ah := window.GetSize(); aw != s.width || ah != s.height {
		log.Fatalf("Window doesn't have the requested size: want %d,%d got %d,%d", s.width, s.height, aw, ah)
	}
	window.MakeContextCurrent()

	// Set window up to be [0, -1.0] -> [width, 1.0], black.
	gl.MatrixMode(gl.MODELVIEW_MATRIX)
	gl.LoadIdentity()
	gl.Translated(-1, 0, 0)
	gl.Scaled(2/float64(s.width), 1.0, 1.0)
	gl.ClearColor(0, 0, 0, 0)

	// Actually start writing data to the buffer.
	s.buffer.GoPushChannel(values, sampleRate)

	for !window.ShouldClose() {
		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			break
		}
		gl.Clear(gl.COLOR_BUFFER_BIT)
		s.drawSignal()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

// drawSignal writes the input wave form(s) out to screen.
func (s *Screen) drawSignal() {
	gl.Begin(gl.LINE_STRIP)
		s.buffer.Each(func(index int, value float64) {
			gl.Vertex2d(float64(index), value)
		})
	gl.End()
}
