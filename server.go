package main

import (
	"fmt"

	"github.com/hypebeast/go-osc/osc"
)

func main() {
    // Default Muse location.
    port := 5000
    server := osc.NewOscServer("127.0.0.1", port)

    fmt.Printf(">> Running OSC server on :%v\n", port)
    server.AddMsgHandler("/osc/address", func(msg *osc.OscMessage) {
        osc.PrintOscMessage(msg)
    })

    server.ListenAndDispatch()
}
