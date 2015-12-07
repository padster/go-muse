gomuse
======

Golang libraries for reading and processing data from a Muse EEG headset.
Requires the [Muse SDK][1] for muse-io.

**Includes**:

 - Running an OSC server to convert the streams into channels.
 - Adds low-level processing of these channels for things like connection status.

**To come**:
 - Simple UIs for displaying realtime data, using go-sound
 - Connection quality detection to avoid needing MuseLab.
 - More complex signal processing.
 - Machine learning integration for feature detection.

Provided is a sample application (server.go) which runs an OSC processor and forwards events to a 
web socket which JS clients can then use to alter a site. To run the example, which increases an image
each time a blink is detected, do the following:

1. Connect your device to the desired OSC port (example uses a device named 'Muse', port = 5000): 
   * `muse-io.exe --device Muse --osc osc.udp://localhost:5000`
2. Run the server (example is using windows):
   * `go build server.go && server.exe`
3. Open [http://localhost:8888/file/viz.html] in a browser.

Note that these can all be done independently (e.g. to reconnect your Muse, to make changes to the server, or to reset the html state) although restarting the server tends to require a browser refresh too.

Additionally, there is currently no detection of bad connection to the Muse device - please use MuseLab's connection strength visualizer for that for now. This can eventually be added to the server too.

Sample application image from the [Tardis Wiki][2].

[1]:https://sites.google.com/a/interaxon.ca/muse-developer-site/museio
[2]:http://tardis.wikia.com/wiki/Weeping_Angel
