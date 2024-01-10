//go:build !windows && !darwin && !cgo
// +build !windows,!darwin,!cgo

package main

var notificationSoundEnabled = false

func playNotificationSound() chan bool { return make(chan bool) }

func initBeep() error { return nil }
