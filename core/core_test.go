package core

import (
	"context"
	"testing"
	"time"

	"gioui.org/x/notify"
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

func TestLoop(t *testing.T) {
	notificationCount := new(int)
	cancellationCount := new(int)
	notifier := mockNotifier{
		cancellationCount: cancellationCount,
		notificationCount: notificationCount,
		t:                 t,
	}

	settings := Settings{
		Duration:  time.Millisecond * 50,
		Frequency: time.Millisecond * 100,
		Sound:     false,
	}

	const timeout = 1000 * time.Millisecond
	// the last notification is unrealiable because of timing
	expectCount := int(timeout/settings.Frequency) - 1
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	loop(ctx, notifier, &settings)

	if *notificationCount < expectCount {
		t.Errorf("Notification count should be at least %d, it was %d", expectCount, *notificationCount)
	}
	if *cancellationCount < expectCount {
		t.Errorf("Cancellation count should be at least %d, it was %d", expectCount, *cancellationCount)
	}
}
