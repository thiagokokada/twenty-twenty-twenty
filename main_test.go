package main

import (
	"testing"

	"gioui.org/x/notify"
)

func TestSendNotification(t *testing.T) {
	notifier, err := notify.NewNotifier()
	if err != nil {
		t.Fatalf("Error while creating a notifier: %v\n", err)
	}
	// ignoring result, because this test does not work in some platforms (e.g.:
	// darwin)
	_ = sendNotification(
		notifier,
		"Test notification title",
		"Test notification text",
		false, // being tested in TestPlayNotificationSound
	)
}
