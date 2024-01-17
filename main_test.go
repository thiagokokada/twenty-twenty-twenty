package main

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"gioui.org/x/notify"
)

var (
	notificationCount int
	cancellationCount int
)

type mockNotifier struct {
	notify.Notifier
	t *testing.T
}

type mockNotification struct {
	*mockNotifier
}

func (notifier mockNotifier) CreateNotification(title, text string) (notify.Notification, error) {
	notificationCount++
	if title != "Time to rest your eyes" {
		notifier.t.Errorf("Title is '%s'", title)
	}
	if text != "Look at 20 feet (~6 meters) away for 0 seconds" {
		notifier.t.Errorf("Text is '%s'", text)
	}
	return &mockNotification{}, nil
}

func (notification mockNotification) Cancel() error {
	cancellationCount++
	return nil
}

func TestTwentyTwentyTwenty(t *testing.T) {
	notificationCount = 0
	cancellationCount = 0
	notifier := mockNotifier{Notifier: nil, t: t}

	duration := new(time.Duration)
	*duration = time.Millisecond * 500

	frequency := new(time.Duration)
	*frequency = time.Second * 1

	notificationSound := new(bool)
	*notificationSound = false

	const timeoutInSec = 5
	// the last notification is unrealiable because of timing
	const expectCount = timeoutInSec - 1
	context, cancel := context.WithTimeout(context.Background(), time.Second*timeoutInSec)

	twentyTwentyTwenty(context, notifier, duration, frequency, notificationSound)
	cancel()

	if notificationCount < expectCount {
		t.Errorf("Notification count should be at least %d, it was %d", expectCount, notificationCount)
	}
	if cancellationCount < expectCount {
		t.Errorf("Cancellation count should be at least %d, it was %d", expectCount, cancellationCount)
	}
}

// The reason those tests exist is to help with development (e.g.: test if
// notification/sound is working). It is useless outside of development purposes
// and needs a proper desktop environment to work, and this is the reason why it
// is not run in CI.
func TestPlayNotificationSound(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}

	err := initNotification()
	if err != nil {
		t.Fatalf("Error while initialising sound: %v\n", err)
	}
	const wait = 10

	log.Println("You should listen to a sound!")
	playSendNotificationSound()
	log.Printf("Waiting %d seconds to ensure that the sound is finished\n", wait)
	time.Sleep(wait * time.Second)

	log.Println("You should listen to another sound!")
	playCancelNotificationSound()
	log.Printf("Waiting %d seconds to ensure that the sound is finished\n", wait)
	time.Sleep(wait * time.Second)
}

func TestSendNotification(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}

	notifier, err := notify.NewNotifier()
	if err != nil {
		t.Fatalf("Error while creating a notifier: %v\n", err)
	}
	log.Println("You should see a notification!")
	notificationSound = new(bool)
	*notificationSound = false
	// ignoring result, because this test does not work in some platforms (e.g.:
	// darwin, because lack of signature)
	_ = sendNotification(
		notifier,
		"Test notification title",
		"Test notification text",
		notificationSound, // being tested in TestPlayNotificationSound
	)
}
