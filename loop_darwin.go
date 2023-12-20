//go:build darwin
// +build darwin

//go:generate go vet ./...
//go:generate go run gioui.org/cmd/gogio -target macos -icon ./eye.png -o TwentyTwentyTwenty.app .
// Signing the code with the adhoc certificate
//go:generate sh -c "codesign -s - TwentyTwentyTwenty.app/*"

package main

import "gioui.org/app"

func loop() {
	app.Main()
}
