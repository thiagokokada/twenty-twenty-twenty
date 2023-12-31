//go:build darwin
// +build darwin

//go:generate go vet ./...
//go:generate sh -c "go run gioui.org/cmd/gogio -ldflags=\"-X 'main.Version=$(git describe --tags --dirty)'\" -o bin/twenty-twenty-twenty -target macos -icon ./eye.png -o TwentyTwentyTwenty.app ."
// Signing the code with the adhoc certificate
//go:generate sh -c "codesign -s - TwentyTwentyTwenty.app/*"

package main

import "gioui.org/app"

func loop() {
	app.Main()
}
