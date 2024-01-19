package notification

import (
	"context"
	"log"
	"time"

	"gioui.org/x/notify"

	snd "github.com/thiagokokada/twenty-twenty-twenty/sound"
)

func Send(
	notifier notify.Notifier,
	title string,
	text string,
	sound *bool,
) notify.Notification {
	if *sound {
		snd.PlaySendNotification()
	}

	notification, err := notifier.CreateNotification(title, text)
	if err != nil {
		log.Printf("Error while sending notification: %v\n", err)
		return nil
	}
	return notification
}

func CancelAfter(
	ctx context.Context,
	notification notify.Notification,
	after *time.Duration,
	sound *bool,
) {
	if notification == nil {
		return
	}

	timer := time.NewTimer(*after)
	select {
	case <-timer.C:
		if *sound {
			snd.PlayCancelNotification()
		}
	case <-ctx.Done(): // avoid playing notification sound if we cancel the context
	}
	err := notification.Cancel()
	if err != nil {
		log.Printf("Error while cancelling notification: %v\n", err)
	}
}
