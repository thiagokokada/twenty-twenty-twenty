package notification

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"gioui.org/x/notify"
)

// The reason this test exist is to help with development (e.g.: test if
// notification). It is useless outside of development purposes and needs a
// proper desktop environment to work, and this is the reason why it is not run
// in CI.
// macOS notes: this does not work in macOS because it needs a signed app bundle
func TestSendWithDuration(t *testing.T) {
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

	go func() {
		time.Sleep(*after)
		log.Println("The notification should have disappeared!")
	}()
	err = SendWithDuration(
		context.Background(),
		notifier,
		after,
		sound,
		"Test notification title",
		"Test notification text",
	)
	if err != nil {
		t.Fatalf("Error while sending notification: %v\n", err)
	}
}
