package core

import (
	"context"
	"strings"
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

func TestParseFlags(t *testing.T) {
	const progname = "twenty-twenty-twenty"
	const version = "test"

	// always return false for sound if disabled
	settings := ParseFlags(progname, []string{}, version, Optional{Sound: false, Systray: false})
	assertEqual(t, settings.Sound, false)

	var tests = []struct {
		args     []string
		settings Settings
	}{
		{[]string{},
			Settings{Duration: time.Second * 20, Frequency: time.Minute * 20, Pause: time.Hour, Sound: true}},
		{[]string{"-duration", "10", "-frequency", "600", "-pause", "1800", "-disable-sound"},
			Settings{Duration: time.Second * 10, Frequency: time.Minute * 10, Pause: time.Minute * 30, Sound: false}},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			settings := ParseFlags(progname, tt.args, version, Optional{Sound: true, Systray: true})
			assertEqual(t, settings, tt.settings)
		})
	}
}

func TestStart(t *testing.T) {
	notifier := newMockNotifier()

	const timeout = time.Second
	// the last notification may or may not come because of timing
	expectCount := int(timeout/testSettings.Frequency) - 1

	Start(notifier, &testSettings, Optional{Sound: true})
	defer Stop()
	time.Sleep(timeout)

	assertGreaterOrEqual(t, *notifier.notificationCount, expectCount)
	assertGreaterOrEqual(t, *notifier.notificationCancelCount, expectCount)
}

func TestPause(t *testing.T) {
	notifier := newMockNotifier()

	const timeout = time.Second
	go func() { time.Sleep(timeout); Stop() }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	callbackPreCalled := false
	callbackPosCalled := false
	Pause(
		ctx, notifier, &testSettings, Optional{},
		func() { callbackPreCalled = true },
		func() { callbackPosCalled = true },
	)

	assertEqual(t, callbackPreCalled, true)
	assertEqual(t, callbackPosCalled, true)
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
	callbackPreCalled := false
	callbackPosCalled := false
	Pause(
		ctx, notifier, &testSettings, Optional{},
		func() { callbackPreCalled = true },
		func() { callbackPosCalled = true },
	)

	assertEqual(t, callbackPreCalled, false)
	assertEqual(t, callbackPosCalled, false)
	assertEqual(t, *notifier.notificationCount, 0)
	assertEqual(t, *notifier.notificationCancelCount, 0)
}
