//go:build darwin
// +build darwin

//go:generate go vet ./...
//go:generate sh -c "go run gioui.org/cmd/gogio -ldflags=\"-X 'main.version=$(git describe --tags --dirty)'\ -s -w" -o bin/twenty-twenty-twenty -target macos -icon ./assets/eye.png -o TwentyTwentyTwenty.app ."
// Signing the code with the adhoc certificate
//go:generate sh -c "codesign -s - TwentyTwentyTwenty.app/*"

package main

import "gioui.org/app"

func loop() {
	app.Main()
}
