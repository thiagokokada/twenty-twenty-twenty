package notification

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync/atomic"
	"time"

	"gioui.org/x/notify"

	snd "github.com/thiagokokada/twenty-twenty-twenty/internal/sound"
)

var notifier atomic.Pointer[notify.Notifier]

func SendWithDuration(
	ctx context.Context,
	duration *time.Duration,
	sound *bool,
	title string,
	text string,
) error {
	notification, err := Send(sound, title, text)
	if err != nil {
		return fmt.Errorf("send notification: %w", err)
	}
	if duration != nil {
		return CancelAfter(ctx, notification, duration, sound)
	}
	return nil
}

func Send(
	sound *bool,
	title string,
	text string,
) (notify.Notification, error) {
	initIfNull()
	if *sound {
		snd.PlaySendNotification(nil)
	}
	return (*notifier.Load()).CreateNotification(title, text)
}

func CancelAfter(
	ctx context.Context,
	notification notify.Notification,
	after *time.Duration,
	sound *bool,
) error {
	timer := time.NewTimer(*after)
	select {
	case <-timer.C:
		if *sound {
			snd.PlayCancelNotification(nil)
		}
	case <-ctx.Done():
		slog.DebugContext(ctx, "Skipping cancel notification sound...")
	}
	slog.DebugContext(ctx, "Cancelling notification...")
	err := notification.Cancel()
	if err != nil {
		return fmt.Errorf("cancel notification: %w", err)
	}
	return nil
}

func SetNotifier(n notify.Notifier) {
	notifier.Store(&n)
}

func initIfNull() {
	if notifier.Load() == nil {
		slog.Debug("Initialising notifier...")
		newNotifier, err := notify.NewNotifier()
		if err != nil {
			log.Fatalf("Error while creating a notifier: %v\n", err)
		}
		swapped := notifier.CompareAndSwap(nil, &newNotifier)
		if !swapped {
			log.Println("Couldn't swap notifier since one is already running")
		}
	}
}
