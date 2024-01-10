//go:build !windows && !darwin && !cgo
// +build !windows,!darwin,!cgo

package main

var notificationSoundEnabled = false

func playNotificationSound() {}

func initBeep() {}
