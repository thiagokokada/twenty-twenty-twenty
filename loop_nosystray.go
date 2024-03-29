//go:build nosystray
// +build nosystray

//go:generate go vet ./...
//go:generate sh -c "go build -tags=nosystray -v -ldflags=\"-X 'main.version=$(git describe --tags --dirty)'\" -o bin/twenty-twenty-twenty"

package main

const systrayEnabled bool = false

func loop() {
	// https://blog.sgmansfield.com/2016/06/how-to-block-forever-in-go/
	select {}
}
