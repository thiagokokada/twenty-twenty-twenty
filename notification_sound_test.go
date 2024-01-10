package main

import (
	"testing"
)

func TestPlayNotificationSound(t *testing.T) {
	err := initBeep()
	if err != nil {
		t.Fatalf("Error while initialising sound: %v\n", err)
	}
	<-playNotificationSound()
	// it takes a while until the sound finishes playing,
	// so uncomment the line below during development
	// time.Sleep(5 * time.Second)
}
