//go:build !windows && !darwin && !cgo
// +build !windows,!darwin,!cgo

package main

const notificationSoundEnabled bool = false

func playNotificationSound() chan bool {
	c := make(chan bool)
	c <- false
	return c
}

func initBeep() error { return nil }
