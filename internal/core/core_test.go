package core

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"gioui.org/x/notify"
	"github.com/thiagokokada/twenty-twenty-twenty/internal/assert"
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

var (
	notifier *mockNotifier
	twenty   = New(
		Features{
			Sound:   false,
			Systray: false,
		},
		Settings{
			Duration:  time.Millisecond * 50,
			Frequency: time.Millisecond * 100,
			Pause:     time.Millisecond * 500,
			Sound:     false,
		},
	)
)

func init() {
	notifier = newMockNotifier()
	notification.Init(notifier)
}

func resetCounters() {
	notifier.notificationCount.Store(0)
	notifier.notificationCancelCount.Store(0)
}

func TestStart(t *testing.T) {
	resetCounters()

	const timeout = time.Second
	// the last notification may or may not come because of timing
	expectCount := int32(timeout/twenty.Settings.Frequency) - 1

	ctx, cancelCtx := context.WithTimeout(context.Background(), timeout)
	defer cancelCtx()
	twenty.Start(ctx)

	assert.GreaterOrEqual(t, notifier.notificationCount.Load(), expectCount)
	assert.GreaterOrEqual(t, notifier.notificationCancelCount.Load(), expectCount)
}

func TestPause(t *testing.T) {
	resetCounters()

	const timeout = time.Second
	ctx, cancelCtx := context.WithTimeout(context.Background(), timeout)
	defer cancelCtx()

	callbackPreCalled := false
	callbackPosCalled := false
	twenty.Pause(
		ctx,
		func() { callbackPreCalled = true },
		func() { callbackPosCalled = true },
	)
	<-ctx.Done()

	assert.Equal(t, callbackPreCalled, true)
	assert.Equal(t, callbackPosCalled, true)
	assert.GreaterOrEqual(t, notifier.notificationCount.Load(), 1)
	assert.GreaterOrEqual(t, notifier.notificationCancelCount.Load(), 1)
}

func TestPauseCancel(t *testing.T) {
	resetCounters()

	ctx, cancelCtx := context.WithCancel(context.Background())
	cancelCtx()

	callbackPreCalled := false
	callbackPosCalled := false
	twenty.Pause(
		ctx,
		func() { callbackPreCalled = true },
		func() { callbackPosCalled = true },
	)

	assert.Equal(t, callbackPreCalled, false)
	assert.Equal(t, callbackPosCalled, false)
	assert.Equal(t, notifier.notificationCount.Load(), 0)
	assert.Equal(t, notifier.notificationCancelCount.Load(), 0)
}

func TestPauseNilCallbacks(t *testing.T) {
	resetCounters()

	const timeout = time.Second
	ctx, cancelCtx := context.WithTimeout(context.Background(), timeout)
	defer cancelCtx()

	twenty.Pause(ctx, nil, nil)

	assert.Equal(t, notifier.notificationCount.Load(), 1)
	assert.Equal(t, notifier.notificationCancelCount.Load(), 1)
}
