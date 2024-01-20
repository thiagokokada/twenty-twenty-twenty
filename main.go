package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"gioui.org/x/notify"

	"github.com/thiagokokada/twenty-twenty-twenty/core"
	ntf "github.com/thiagokokada/twenty-twenty-twenty/notification"
	"github.com/thiagokokada/twenty-twenty-twenty/sound"
)

var (
	version  = "development"
	notifier notify.Notifier
	optional core.Optional
	settings core.Settings
)

func main() {
	optional = core.Optional{Sound: sound.Enabled, Systray: systrayEnabled}
	settings = core.ParseFlags(os.Args[0], os.Args[1:], version, optional)
	var err error

	// only init Beep if notification sound is enabled, otherwise we will cause
	// unnecessary noise in the speakers (and also increased memory usage)
	if settings.Sound {
		err = sound.Init()
		if err != nil {
			log.Fatalf("Error while initialising sound: %v\n", err)
		}
	}

	notifier, err = notify.NewNotifier()
	if err != nil {
		log.Fatalf("Error while creating a notifier: %v\n", err)
	}

	err = ntf.Send(
		context.Background(),
		notifier,
		&settings.Duration,
		&settings.Sound,
		"Starting 20-20-20",
		fmt.Sprintf("You will see a notification every %.f minutes(s)", settings.Frequency.Minutes()),
	)
	if err != nil {
		log.Fatalf("Test notification failed: %v. Exiting...", err)
	}

	go core.Start(notifier, &settings, optional)
	loop()
}
