//go:build !windows && !darwin && !cgo
// +build !windows,!darwin,!cgo

package main

const notificationSoundEnabled bool = false

func playNotificationSound1() {
	panic("Not implemented")
}

func playNotificationSound2() {
	panic("Not implemented")
}

func initNotification() error {
	panic("Not implemented")
}
