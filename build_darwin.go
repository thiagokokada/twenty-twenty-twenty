//go:build darwin
// +build darwin

//go:generate go run gioui.org/cmd/gogio -target macos -icon ./eye.png -o TwentyTwentyTwenty.app .
//go:generate sh -c "codesign -s - TwentyTwentyTwenty.app/*"

package main
