package core

import (
	"context"
	"testing"
	"time"

	"gioui.org/x/notify"
)

type mockNotifier struct {
	notify.Notifier
	notificationCancelCount *int
	notificationCount       *int
}

type mockNotification struct {
	*mockNotifier
}

func (n mockNotifier) CreateNotification(title, text string) (notify.Notification, error) {
	*n.notificationCount++
	return &mockNotification{mockNotifier: &n}, nil
}

func (n mockNotification) Cancel() error {
	*n.mockNotifier.notificationCancelCount++
	return nil
}

func newMockNotifier() mockNotifier {
	return mockNotifier{
		notificationCancelCount: new(int),
		notificationCount:       new(int),
	}
}

var testSettings = Settings{
	Duration:  time.Millisecond * 50,
	Frequency: time.Millisecond * 100,
	Pause:     time.Millisecond * 500,
	Sound:     false,
}

func TestStartAndStop(t *testing.T) {
	notifier := newMockNotifier()

	const timeout = 1000 * time.Millisecond
	// the last notification is unrealiable because of timing
	expectCount := int(timeout/testSettings.Frequency) - 1

	go func() { time.Sleep(timeout); Stop() }()
	Start(notifier, &testSettings)

	if *notifier.notificationCount < expectCount {
		t.Errorf(
			"Notification count should be at least %d, it was %d",
			expectCount,
			notifier.notificationCount,
		)
	}
	if *notifier.notificationCancelCount < expectCount {
		t.Errorf(
			"Cancellation count should be at least %d, it was %d",
			expectCount,
			notifier.notificationCancelCount,
		)
	}
}

func TestPause(t *testing.T) {
	notifier := newMockNotifier()

	const timeout = 1000 * time.Millisecond
	// the last notification is unrealiable because of timing
	expectCount := int((testSettings.Pause-timeout)/testSettings.Frequency) - 1
	go func() { time.Sleep(timeout); Stop() }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	callbackCalled := false
	Pause(ctx, notifier, &testSettings, func() { callbackCalled = true })

	if !callbackCalled {
		t.Error("Callback should have been called")
	}
	if *notifier.notificationCount < expectCount {
		t.Errorf(
			"Notification count should be at least %d, it was %d",
			expectCount,
			notifier.notificationCount,
		)
	}
	if *notifier.notificationCancelCount < expectCount {
		t.Errorf(
			"Cancellation count should be at least %d, it was %d",
			expectCount,
			notifier.notificationCancelCount,
		)
	}
}

func TestCancelledPause(t *testing.T) {
	notifier := newMockNotifier()

	const timeout = 1000 * time.Millisecond
	expectCount := 0
	go func() { time.Sleep(timeout); Stop() }()

	// will be cancelled before the timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout/10)
	defer cancel()
	callbackCalled := false
	Pause(ctx, notifier, &testSettings, func() { callbackCalled = true })

	if callbackCalled {
		t.Error("Callback should not have been called")
	}
	if *notifier.notificationCount < expectCount {
		t.Errorf(
			"Notification count should be at least %d, it was %d",
			expectCount,
			notifier.notificationCount,
		)
	}
	if *notifier.notificationCancelCount < expectCount {
		t.Errorf(
			"Cancellation count should be at least %d, it was %d",
			expectCount,
			notifier.notificationCancelCount,
		)
	}
}
