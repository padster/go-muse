
package main

// To get the muse device to stream to port 5000:
// muse-io.exe --device Muse --osc osc.udp://localhost:5000

import (
	"fmt"

	"github.com/hypebeast/go-osc/osc"

	"golang.org/x/net/websocket"
	"net/http"
	// "runtime"
	"time"
)

const (
	// From testing with MuseLab, eeg[0] only drops under about 700 when I blink.
	// NOTE: can't use /muse/elements/blink, it's too laggy.
	BLINK_EEG_THRESHOLD = 700

	// Avoid double-counting blinks by forcing a time gap between them.
	BLINK_TIME_THRESHOLD_MS = 150.0
)


// A server that forwards OSC events to clients over websockets.
type Server struct {
	activeClients []*websocket.Conn
	activeMessages []chan string
}

// Utility to serve the static files used for the app.
func serveStaticFiles(fromDirectory string, toHttpPrefix string) {
	asPath := fmt.Sprintf("/%s/", toHttpPrefix)
	fmt.Printf("Serving %s as %s\n", fromDirectory, toHttpPrefix)
	fs := http.FileServer(http.Dir(fromDirectory))
	http.Handle(asPath, disableCache(http.StripPrefix(asPath, fs)))
}
func disableCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		h.ServeHTTP(w, r)
	})
}

// Save the connection, replacing any previous one.
func (s *Server) socketHandler(ws *websocket.Conn) {
	fmt.Printf("Client connected!\n")

	// PICK: Properly clean up on client disconnect instead?
	if len(s.activeClients) == 1 {
		s.activeClients[0].Close()
		close(s.activeMessages[0])
		s.activeClients[0] = ws
		s.activeMessages[0] = make(chan string)
	} else {
		s.activeClients = append(s.activeClients, ws)
		s.activeMessages = append(s.activeMessages, make(chan string))
	}

	// Continually wait for messages to forward to client.
	for msg := range s.activeMessages[0] {
		ws.Write([]byte(msg))
	}
}

// Launches a server to serve the static files and respond to websocket connections.
func (s *Server) openWebServer(port int) {
	fmt.Printf("Opening WebServer on :%d...\n", port)

	serveStaticFiles("./static", "file")
	http.Handle("/sock", websocket.Handler(s.socketHandler))
	
	// NOTE: blocking.
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}

// Launches a server to stream OSC data, find events, and forward to clients.
func (s *Server) openOscServer(port int) {
	fmt.Printf("Opening OSCserver on :%d...\n", port)
	server := &osc.Server{
		Addr: fmt.Sprintf("127.0.0.1:%d", port),
	}

	lastBlinkTime := time.Now()
	server.Handle("/muse/eeg", func(msg *osc.Message) {
		maybeBlink := (msg.Arguments[0].(float32) < BLINK_EEG_THRESHOLD)
		if maybeBlink {
			blinkTime := time.Now()
			msSinceLast := blinkTime.Sub(lastBlinkTime).Seconds() * 1000

			if msSinceLast > BLINK_TIME_THRESHOLD_MS {
				lastBlinkTime = blinkTime
				s.HandleBlink()
			}
		}
	})
	server.ListenAndServe()
}

// Process a single blink identified by the server.
func (s *Server) HandleBlink() {
	fmt.Printf("Blinked!\n")
	msg := fmt.Sprintf("%d", 0) // Blink = '0'
	for _, client := range s.activeMessages {
		client <- msg
	}
}

func main() {
	server := &Server {
		make([]*websocket.Conn, 0, 1),
		make([]chan string, 0, 1),
	}

	go server.openOscServer(5000)
	server.openWebServer(8888)
}
