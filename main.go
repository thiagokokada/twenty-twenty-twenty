package main

import (
	"context"
	"fmt"
	"log"

	"gioui.org/x/notify"

	"github.com/thiagokokada/twenty-twenty-twenty/core"
	n "github.com/thiagokokada/twenty-twenty-twenty/notification"
	snd "github.com/thiagokokada/twenty-twenty-twenty/sound"
)

var (
	version  = "development"
	notifier notify.Notifier
	settings core.Settings
)

func main() {
	settings = core.ParseFlags(version, systrayEnabled, snd.Enabled)
	var err error

	// only init Beep if notification sound is enabled, otherwise we will cause
	// unnecessary noise in the speakers (and also increased memory usage)
	if settings.Sound {
		err = snd.Init()
		if err != nil {
			log.Fatalf("Error while initialising sound: %v\n", err)
		}
	}

	notifier, err = notify.NewNotifier()
	if err != nil {
		log.Fatalf("Error while creating a notifier: %v\n", err)
	}

	notification := n.Send(
		notifier,
		"Starting 20-20-20",
		fmt.Sprintf("You will see a notification every %.f minutes(s)", settings.Frequency.Minutes()),
		&settings.Sound,
	)
	if notification == nil {
		log.Fatalf("Test notification failed, exiting...")
	}
	go n.CancelAfter(context.Background(), notification, &settings.Duration, &settings.Sound)

	core.Start(notifier, &settings)
	loop()
}
