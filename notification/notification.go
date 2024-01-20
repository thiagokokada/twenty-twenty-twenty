package notification

import (
	"context"
	"fmt"
	"time"

	"gioui.org/x/notify"

	snd "github.com/thiagokokada/twenty-twenty-twenty/sound"
)

func SendWithDuration(
	ctx context.Context,
	notifier notify.Notifier,
	duration *time.Duration,
	sound *bool,
	title string,
	text string,
) error {
	notification, err := Send(notifier, sound, title, text)
	if err != nil {
		return fmt.Errorf("send notification: %w", err)
	}
	if duration != nil {
		return CancelAfter(ctx, notification, duration, sound)
	}
	return nil
}

func Send(
	notifier notify.Notifier,
	sound *bool,
	title string,
	text string,
) (notify.Notification, error) {
	if *sound {
		snd.PlaySendNotification(sndCallback)
	}
	return notifier.CreateNotification(title, text)
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
			snd.PlayCancelNotification(sndCallback)
		}
	case <-ctx.Done(): // avoid playing notification sound if we cancel the context
	}
	err := notification.Cancel()
	if err != nil {
		return fmt.Errorf("cancel notification: %w", err)
	}
	return nil
}

func sndCallback() {}
