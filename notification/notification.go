package notification

import (
	"context"
	"fmt"
	"log"
	"time"

	"gioui.org/x/notify"

	snd "github.com/thiagokokada/twenty-twenty-twenty/sound"
)

var notifier notify.Notifier

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
		snd.PlaySendNotification(func() {})
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
			snd.PlayCancelNotification(func() {})
		}
	case <-ctx.Done(): // avoid playing notification sound if we cancel the context
	}
	err := notification.Cancel()
	if err != nil {
		return fmt.Errorf("cancel notification: %w", err)
	}
	return nil
}

func SetNotifier(n notify.Notifier) {
	notifier = n
}

func initIfNull() {
	if notifier == nil {
		notifier, err := notify.NewNotifier()
		if err != nil {
			log.Fatalf("Error while creating a notifier: %v\n", err)
		}
		SetNotifier(notifier)
	}
}
