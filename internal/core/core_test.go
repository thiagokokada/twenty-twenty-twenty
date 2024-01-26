package core

import (
	"cmp"
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"gioui.org/x/notify"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/notification"
)

type mockNotifier struct {
	notify.Notifier
	notificationCancelCount atomic.Int32
	notificationCount       atomic.Int32
}

type mockNotification struct {
	*mockNotifier
}

func (n *mockNotifier) CreateNotification(title, text string) (notify.Notification, error) {
	n.notificationCount.Add(1)
	return &mockNotification{mockNotifier: n}, nil
}

func (n *mockNotification) Cancel() error {
	n.notificationCancelCount.Add(1)
	return nil
}

func newMockNotifier() *mockNotifier {
	return &mockNotifier{
		notificationCancelCount: atomic.Int32{},
		notificationCount:       atomic.Int32{},
	}
}

func assertEqual[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("got: %v; want: %v", actual, expected)
	}
}

func assertGreaterOrEqual[T cmp.Ordered](t *testing.T, actual, expected T) {
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
		{[]string{"-duration", "10", "-frequency", "600", "-pause", "1800", "-disable-sound", "-verbose"},
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
	notification.SetNotifier(notifier)

	const timeout = time.Second
	// the last notification may or may not come because of timing
	expectCount := int32(timeout/testSettings.Frequency) - 1

	Start(&testSettings, Optional{Sound: true})
	defer Stop()
	time.Sleep(timeout)

	assertGreaterOrEqual(t, notifier.notificationCount.Load(), expectCount)
	assertGreaterOrEqual(t, notifier.notificationCancelCount.Load(), expectCount)
}

func TestPause(t *testing.T) {
	notifier := newMockNotifier()
	notification.SetNotifier(notifier)

	const timeout = time.Second
	go func() { time.Sleep(timeout); Stop() }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	callbackPreCalled := false
	callbackPosCalled := false
	Pause(
		ctx, &testSettings, Optional{},
		func() { callbackPreCalled = true },
		func() { callbackPosCalled = true },
	)
	<-Ctx().Done()

	assertEqual(t, callbackPreCalled, true)
	assertEqual(t, callbackPosCalled, true)
	assertGreaterOrEqual(t, notifier.notificationCount.Load(), 1)
	assertGreaterOrEqual(t, notifier.notificationCancelCount.Load(), 1)
}

func TestPauseCancel(t *testing.T) {
	notifier := newMockNotifier()
	notification.SetNotifier(notifier)

	const timeout = time.Second
	go func() { time.Sleep(timeout); Stop() }()

	// will be cancelled before the timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout/10)
	defer cancel()
	callbackPreCalled := false
	callbackPosCalled := false
	Pause(
		ctx, &testSettings, Optional{},
		func() { callbackPreCalled = true },
		func() { callbackPosCalled = true },
	)
	<-Ctx().Done()

	assertEqual(t, callbackPreCalled, false)
	assertEqual(t, callbackPosCalled, false)
	assertEqual(t, notifier.notificationCount.Load(), 0)
	assertEqual(t, notifier.notificationCancelCount.Load(), 0)
}

func TestPauseNilCallbacks(t *testing.T) {
	notifier := newMockNotifier()
	notification.SetNotifier(notifier)

	const timeout = time.Second
	go func() { time.Sleep(timeout); Stop() }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	Pause(ctx, &testSettings, Optional{}, nil, nil)

	assertEqual(t, notifier.notificationCount.Load(), 1)
	assertEqual(t, notifier.notificationCancelCount.Load(), 1)
}
