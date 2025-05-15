module github.com/ClifHouck/unified/examples/doorbell

go 1.24.3

require (
	github.com/ClifHouck/unified v0.1.0
	github.com/gopxl/beep v1.4.1
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.9.1
)

replace github.com/ClifHouck/unified => ../../

require (
	github.com/coder/websocket v1.8.13 // indirect
	github.com/ebitengine/oto/v3 v3.1.0 // indirect
	github.com/ebitengine/purego v0.7.1 // indirect
	github.com/hajimehoshi/go-mp3 v0.3.4 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	golang.org/x/sys v0.29.0 // indirect
)
