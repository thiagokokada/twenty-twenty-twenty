package notification

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"gioui.org/x/notify"
)

// The reason those tests exist is to help with development (e.g.: test if
// notification/sound is working). It is useless outside of development purposes
// and needs a proper desktop environment to work, and this is the reason why it
// is not run in CI.
func TestSendNotification(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}

	notifier, err := notify.NewNotifier()
	if err != nil {
		t.Fatalf("Error while creating a notifier: %v\n", err)
	}
	log.Println("You should see a notification!")

	sound := new(bool)
	*sound = false // being tested in sound package
	after := new(time.Duration)
	*after = time.Duration(5) * time.Second

	notification := Send(notifier, "Test notification title", "Test notification text", sound)
	CancelAfter(context.Background(), notification, after, sound)
}
