// Demo usage of ui.NewScreen OpenGL rendering. 
package main

import (
	"math"
	"time"

	"github.com/padster/go-muse/ui"
)

func main() {
	// Generate a sum of three sin waves, 200Hz...
	wave := make(chan float64)
	go func() {
		x := 0.0
		for {
			x += 0.01
			wave <- (math.Sin(x) + math.Sin(2 * x) + math.Sin(4 * x)) / 3.0
			time.Sleep(5 * time.Millisecond)
		}
	}()

	// ...and draw them to screen. 
	s := ui.NewScreen(1200, 600)
	s.Render(wave, 3 /* sampleRate */)
}
