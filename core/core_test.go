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

func assertEqual[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func assertGreaterOrEqual(t *testing.T, actual, expected int) {
	t.Helper()
	if actual < expected {
		t.Errorf("got: %v; want: >=%v", actual, expected)
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
	// the last notification may or may not come because of timing
	expectCount := int(timeout/testSettings.Frequency) - 1

	go func() { time.Sleep(timeout); Stop() }()
	Start(notifier, &testSettings)

	assertGreaterOrEqual(t, *notifier.notificationCount, expectCount)
	assertGreaterOrEqual(t, *notifier.notificationCancelCount, expectCount)
}

func TestPause(t *testing.T) {
	notifier := newMockNotifier()

	const timeout = time.Second
	go func() { time.Sleep(timeout); Stop() }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	callbackCalled := false
	Pause(ctx, notifier, &testSettings, func() { callbackCalled = true })

	assertEqual(t, callbackCalled, true)
	assertGreaterOrEqual(t, *notifier.notificationCount, 1)
	assertGreaterOrEqual(t, *notifier.notificationCancelCount, 1)
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

	assertEqual(t, callbackCalled, false)
	assertGreaterOrEqual(t, *notifier.notificationCount, 0)
	assertGreaterOrEqual(t, *notifier.notificationCancelCount, 0)
}
