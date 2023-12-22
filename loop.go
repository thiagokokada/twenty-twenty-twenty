//go:build !darwin
// +build !darwin

//go:generate go vet ./...
//go:generate sh -c "go build -v -ldflags=\"-X 'main.Version=$(git describe --tags --dirty)'\" -o bin/twenty-twenty-twenty"

package main

func loop() {
	// https://blog.sgmansfield.com/2016/06/how-to-block-forever-in-go/
	select {}
}
