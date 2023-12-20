//go:build !darwin
// +build !darwin

//go:generate go vet ./...
//go:generate go build -o bin/twenty-twenty-twenty

package main

func mainLoop() {
	select {}
}
