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

func assertAtLeast(t *testing.T, actual int, expected int) {
	t.Helper()
	if actual < expected {
		t.Errorf("got: %v; want: %v", actual, expected)
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

	const timeout = time.Second
	// the last notification may or may not coming because of timing
	expectCount := int(timeout/testSettings.Frequency) - 1

	go func() { time.Sleep(timeout); Stop() }()
	Start(notifier, &testSettings)

	assertAtLeast(t, *notifier.notificationCount, expectCount)
	assertAtLeast(t, *notifier.notificationCancelCount, expectCount)
}

func TestPause(t *testing.T) {
	notifier := newMockNotifier()

	const timeout = time.Second
	go func() { time.Sleep(timeout); Stop() }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	callbackCalled := false
	Pause(ctx, notifier, &testSettings, func() { callbackCalled = true })

	if !callbackCalled {
		t.Error("Callback should have been called")
	}
	assertAtLeast(t, *notifier.notificationCount, 1)
	assertAtLeast(t, *notifier.notificationCancelCount, 1)
}

func TestPauseCancel(t *testing.T) {
	notifier := newMockNotifier()

	const timeout = time.Second
	go func() { time.Sleep(timeout); Stop() }()

	// will be cancelled before the timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout/10)
	defer cancel()
	callbackCalled := false
	Pause(ctx, notifier, &testSettings, func() { callbackCalled = true })

	if callbackCalled {
		t.Error("Callback should not have been called")
	}
	assertAtLeast(t, *notifier.notificationCount, 0)
	assertAtLeast(t, *notifier.notificationCancelCount, 0)
}
