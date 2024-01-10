package main

import (
	"log"
	"testing"
	"time"

	"gioui.org/x/notify"
)

// The reason those tests exist is to help with development (e.g.: test if
// notification/sound is working). It is useless outside of development purposes
// and needs a proper desktop environment to work, and this is the reason why it
// is not run in CI.

func TestPlayNotificationSound(t *testing.T) {
	err := initBeep()
	if err != nil {
		t.Fatalf("Error while initialising sound: %v\n", err)
	}
	log.Println("You should listen to a sound!")
	<-playNotificationSound()
	log.Println("Waiting 5 seconds to ensure that the sound is finished")
	time.Sleep(5 * time.Second)
}

func TestSendNotification(t *testing.T) {
	notifier, err := notify.NewNotifier()
	if err != nil {
		t.Fatalf("Error while creating a notifier: %v\n", err)
	}
	log.Println("You should see a notification!")
	// ignoring result, because this test does not work in some platforms (e.g.:
	// darwin, because lack of signature)
	_ = sendNotification(
		notifier,
		"Test notification title",
		"Test notification text",
		false, // being tested in TestPlayNotificationSound
	)
}
