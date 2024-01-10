package main

import (
	"testing"

	"gioui.org/x/notify"
)

// The reason this test exist is to help with development (e.g.: test if
// notification is working). It is useless outside of development purposes and
// needs a proper desktop environment to work, and this is the reason why it is
// not run in CI.
func TestSendNotification(t *testing.T) {
	notifier, err := notify.NewNotifier()
	if err != nil {
		t.Fatalf("Error while creating a notifier: %v\n", err)
	}
	// ignoring result, because this test does not work in some platforms (e.g.:
	// darwin, because lack of signature)
	_ = sendNotification(
		notifier,
		"Test notification title",
		"Test notification text",
		false, // being tested in TestPlayNotificationSound
	)
}
