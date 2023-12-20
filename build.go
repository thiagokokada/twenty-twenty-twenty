//go:build !darwin
// +build !darwin

//go:generate go vet -tags=novulkan,nowayland,nox11 ./...
//go:generate go build -tags=novulkan,nowayland,nox11 -o bin/twenty-twenty-twenty

package main
