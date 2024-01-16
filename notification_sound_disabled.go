//go:build !windows && !darwin && !cgo
// +build !windows,!darwin,!cgo

package main

const notificationSoundEnabled bool = false

func playSendNotificationSound() {
	panic("Not implemented")
}

func playCancelNotificationSound() {
	panic("Not implemented")
}

func initNotification() error {
	panic("Not implemented")
}
