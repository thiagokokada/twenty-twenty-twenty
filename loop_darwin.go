//go:build darwin
// +build darwin

//go:generate go vet ./...
//go:generate sh -c "go run gioui.org/cmd/gogio -arch=$(uname -m | sed 's|x86_64|amd64|') -target=macos -ldflags=\"-X 'main.version=$(git describe --tags --dirty)' -s -w\" -icon=./assets/eye.png -o=bin/TwentyTwentyTwenty.app ."
// Signing the code with the adhoc certificate
//go:generate sh -c "codesign -s - bin/TwentyTwentyTwenty.app"

package main

import (
	"fyne.io/systray"
)

func loop() {
	systray.Run(onReady, onExit)
}
