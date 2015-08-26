// muse-io.exe --device Muse --osc osc.udp://localhost:5000
package main

import (
	"fmt"

	"github.com/hypebeast/go-osc/osc"

	// Hack - show an image
    "image"
    "image/color"
    "image/draw"

	"github.com/google/gxui"
    // "github.com/google/gxui/drivers/gl"
    "github.com/google/gxui/themes/dark"
)

func main() {
	// Default Muse location.
	port := 5000

	server := &osc.Server{
		Addr: "127.0.0.1:5000",
	}

	fmt.Printf(">> Running OSC server on :%v\n", port)
	lastWasBlink := false
	server.Handle("/muse/elements/blink", func(msg *osc.Message) {
		blink := (msg.Arguments[0].(int32) == 1)
		if blink && !lastWasBlink {
			fmt.Printf("Blinked!\n")
		}
		lastWasBlink = blink
	})

	server.ListenAndServe()
    // gl.StartDriver(appMain)
}

func appMain(driver gxui.Driver) {
    width, height := 640, 480
    m := image.NewRGBA(image.Rect(0, 0, width, height))
    blue := color.RGBA{0, 0, 255, 255}
    draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)

    // The themes create the content. Currently only a dark theme is offered for GUI elements.
    theme := dark.CreateTheme(driver)
    img := theme.CreateImage()
    window := theme.CreateWindow(width, height, "Image viewer")
    texture := driver.CreateTexture(m, 1.0)
    img.SetTexture(texture)
    window.AddChild(img)
    window.OnClose(driver.Terminate)
}

