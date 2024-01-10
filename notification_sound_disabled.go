//go:build linux && !cgo
// +build linux,!cgo

package main

var notificationSoundEnabled = false

func playNotificationSound() {}

func initBeep() {}
