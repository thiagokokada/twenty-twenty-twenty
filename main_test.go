package main

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"gioui.org/x/notify"

	s "github.com/thiagokokada/twenty-twenty-twenty/settings"
)

type mockNotifier struct {
	notify.Notifier
	cancellationCount *int
	notificationCount *int
	t                 *testing.T
}

type mockNotification struct {
	*mockNotifier
}

func (n mockNotifier) CreateNotification(title, text string) (notify.Notification, error) {
	*n.notificationCount++
	if title != "Time to rest your eyes" {
		n.t.Errorf("Title is '%s'", title)
	}
	if text != "Look at 20 feet (~6 meters) away for 0 seconds" {
		n.t.Errorf("Text is '%s'", text)
	}
	return &mockNotification{mockNotifier: &n}, nil
}

func (n mockNotification) Cancel() error {
	*n.mockNotifier.cancellationCount++
	return nil
}

func TestTwentyTwentyTwenty(t *testing.T) {
	notificationCount := new(int)
	cancellationCount := new(int)
	notifier := mockNotifier{
		cancellationCount: cancellationCount,
		notificationCount: notificationCount,
		t:                 t,
	}

	settings := s.Settings{
		Duration:  time.Millisecond * 50,
		Frequency: time.Millisecond * 100,
		Sound:     false,
	}

	const timeout = 1000 * time.Millisecond
	// the last notification is unrealiable because of timing
	expectCount := int(timeout/settings.Frequency) - 1
	ctx, ctxCancel := context.WithTimeout(context.Background(), timeout)

	twentyTwentyTwenty(ctx, notifier, &settings)
	ctxCancel()

	if *notificationCount < expectCount {
		t.Errorf("Notification count should be at least %d, it was %d", expectCount, *notificationCount)
	}
	if *cancellationCount < expectCount {
		t.Errorf("Cancellation count should be at least %d, it was %d", expectCount, *cancellationCount)
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

	err := initSound()
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

	sound := new(bool)
	*sound = false
	// ignoring result, because this test does not work in some platforms (e.g.:
	// darwin, because lack of signature)
	_ = sendNotification(
		notifier,
		"Test notification title",
		"Test notification text",
		sound, // being tested in TestPlayNotificationSound
	)
}
