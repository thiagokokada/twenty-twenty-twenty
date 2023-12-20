//go:build !darwin
// +build !darwin

//go:generate go vet ./...
//go:generate go run honnef.co/go/tools/cmd/staticcheck ./...
//go:generate go build -o bin/twenty-twenty-twenty

package main
