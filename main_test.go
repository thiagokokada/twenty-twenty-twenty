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
	err := initNotification()
	if err != nil {
		t.Fatalf("Error while initialising sound: %v\n", err)
	}
	const wait = 10

	log.Println("You should listen to a sound!")
	playNotificationSound1()
	log.Printf("Waiting %d seconds to ensure that the sound is finished", wait)
	time.Sleep(wait * time.Second)

	log.Println("You should listen to another sound!")
	playNotificationSound2()
	log.Printf("Waiting %d seconds to ensure that the sound is finished", wait)
	time.Sleep(wait * time.Second)
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
