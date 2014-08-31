// Demo usage of ui.NewScreen OpenGL rendering. 
package main

import (
	"github.com/padster/go-muse/ui"
)

func main() {
	// TODO: Pipe in my own float64 channel.
	s := ui.NewScreen(1000, 600)
	s.Start()
}