package sound

import (
	"os"
	"testing"
	"time"
)

// The reason this test exist is to help with development (e.g.: test if sound
// is working). It is useless outside of development purposes and needs a proper
// desktop environment to work, and this is the reason why it is not run in CI.
func TestPlaySendAndCancelNotification(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}

	err := Init(true)
	if err != nil {
		t.Fatalf("Error while initialising sound: %v\n", err)
	}
	// this shouldn't cut the sound during playback
	suspendDone := make(chan bool)
	go func() {
		SuspendAfter(time.Second * 10)
		suspendDone <- true
	}()

	soundDone := make(chan bool)
	t.Log("You should listen to a sound!")
	PlaySendNotification(func() { soundDone <- true })
	<-soundDone

	t.Log("You should listen to another sound!")
	PlayCancelNotification(func() { soundDone <- true })
	<-soundDone

	<-suspendDone
}
