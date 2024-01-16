//go:build !nosystray && !darwin
// +build !nosystray,!darwin

//go:generate go vet ./...
//go:generate sh -c "go build -v -ldflags=\"-X 'main.version=$(git describe --tags --dirty)'\" -o bin/twenty-twenty-twenty"

package main

import "fyne.io/systray"

func loop() {
	systray.Run(onReady, onExit)
}
