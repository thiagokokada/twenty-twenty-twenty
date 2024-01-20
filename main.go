package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/x/notify"

	"github.com/thiagokokada/twenty-twenty-twenty/core"
	"github.com/thiagokokada/twenty-twenty-twenty/notification"
	"github.com/thiagokokada/twenty-twenty-twenty/sound"
)

var (
	version  = "development"
	optional core.Optional
	settings core.Settings
)

func main() {
	optional = core.Optional{Sound: sound.Enabled, Systray: systrayEnabled}
	settings = core.ParseFlags(os.Args[0], os.Args[1:], version, optional)

	// only init Beep if notification sound is enabled, otherwise we will cause
	// unnecessary noise in the speakers (and also increased memory usage)
	if settings.Sound {
		err := sound.Init()
		if err != nil {
			log.Fatalf("Error while initialising sound: %v\n", err)
		}
	}

	notifier, err := notify.NewNotifier()
	if err != nil {
		log.Fatalf("Error while creating a notifier: %v\n", err)
	}
	notification.Init(notifier)

	sentNotification, err := notification.Send(
		&settings.Sound,
		"Starting 20-20-20",
		fmt.Sprintf("You will see a notification every %.f minutes(s)", settings.Frequency.Minutes()),
	)
	if err != nil {
		log.Fatalf("Test notification failed: %v. Exiting...", err)
	}
	// we need to start notification cancellation in a goroutine to show the
	// systray as soon as possible (since it depends on the loop() call), but we
	// also need to give it access to the core.Ctx to cancel it if necessary
	core.Start(&settings, optional)
	go notification.CancelAfter(core.Ctx(), sentNotification, &settings.Duration, &settings.Sound)

	loop()
}
