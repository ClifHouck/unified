module github.com/ClifHouck/unified/examples/doorbell

go 1.23.8

require (
	github.com/ClifHouck/unified v0.0.0-20250504131910-33720fcc5f07
	github.com/gopxl/beep v1.4.1
	github.com/sirupsen/logrus v1.9.3
)

replace github.com/ClifHouck/unified => ../../

require (
	github.com/coder/websocket v1.8.13 // indirect
	github.com/ebitengine/oto/v3 v3.1.0 // indirect
	github.com/ebitengine/purego v0.7.1 // indirect
	github.com/hajimehoshi/go-mp3 v0.3.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
