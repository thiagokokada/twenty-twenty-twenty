package sound

import (
	"log"
	"os"
	"testing"
	"time"
)

// The reason those tests exist is to help with development (e.g.: test if
// notification/sound is working). It is useless outside of development purposes
// and needs a proper desktop environment to work, and this is the reason why it
// is not run in CI.
func TestPlayNotificationSound(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}

	err := Init()
	if err != nil {
		t.Fatalf("Error while initialising sound: %v\n", err)
	}
	const wait = 10

	log.Println("You should listen to a sound!")
	PlaySendNotification()
	log.Printf("Waiting %d seconds to ensure that the sound is finished\n", wait)
	time.Sleep(wait * time.Second)

	log.Println("You should listen to another sound!")
	PlayCancelNotification()
	log.Printf("Waiting %d seconds to ensure that the sound is finished\n", wait)
	time.Sleep(wait * time.Second)
}
