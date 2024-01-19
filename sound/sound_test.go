package sound

import (
	"log"
	"os"
	"testing"
)

// The reason this test exist is to help with development (e.g.: test if sound
// is working). It is useless outside of development purposes and needs a proper
// desktop environment to work, and this is the reason why it is not run in CI.
func TestPlaySendAndCancelNotification(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}

	err := Init()
	if err != nil {
		t.Fatalf("Error while initialising sound: %v\n", err)
	}
	const wait = 10

	done := make(chan bool)
	log.Println("You should listen to a sound!")
	PlaySendNotification(func() { done <- true })
	<-done

	log.Println("You should listen to another sound!")
	PlayCancelNotification(func() { done <- true })
	<-done
}
