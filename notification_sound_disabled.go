//go:build !windows && !darwin && !cgo
// +build !windows,!darwin,!cgo

package main

const notificationSoundEnabled bool = false

func playNotificationSound1() {}

func playNotificationSound2() {}

func initNotification() error { return nil }
