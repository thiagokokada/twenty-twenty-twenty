//go:build !nosystray
// +build !nosystray

package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"time"

	"fyne.io/systray"
)

const systrayEnabled bool = true

//go:embed assets/eye_light.ico
var data []byte

type menuItems struct {
	mEnabled *systray.MenuItem
	mPause   *systray.MenuItem
	mQuit    *systray.MenuItem
	mSound   *systray.MenuItem
}

func resumeTwentyTwentyTwentyAfter(
	ctx context.Context,
	ctxCancel context.CancelFunc,
	settings *appSettings,
	menu *menuItems,
) {
	log.Printf("Pausing twenty-twenty-twenty for %.f hour...\n", settings.pause.Hours())
	mainCtxCancel() // cancelling current twenty-twenty-twenty goroutine
	timer := time.NewTimer(settings.pause)
	cancelCtx, cancelCtxCancel := context.WithCancel(context.Background())

	select {
	case <-timer.C:
		notification := sendNotification(
			notifier,
			"Resuming 20-20-20",
			fmt.Sprintf("You will see a notification every %.f minutes(s)", settings.frequency.Minutes()),
			&settings.sound,
		)
		if notification == nil {
			log.Printf("Resume notification failed...")
		}
		go cancelNotificationAfter(cancelCtx, notification, settings)
		runTwentyTwentyTwenty(notifier, settings)

		menu.mEnabled.Enable()
		menu.mPause.Uncheck()
	case <-ctx.Done():
	}
	cancelCtxCancel()
	ctxCancel() // make sure the current context is closed
}

func onReady() {
	systray.SetIcon(data)
	systray.SetTooltip("TwentyTwentyTwenty")
	mEnabled := systray.AddMenuItemCheckbox("Enabled", "Enable twenty-twenty-twenty", true)
	mPause := systray.AddMenuItemCheckbox(
		fmt.Sprintf("Pause for %.f hour", settings.pause.Hours()),
		fmt.Sprintf("Pause twenty-twenty-twenty for %.f hour", settings.pause.Hours()),
		false,
	)
	mSound := new(systray.MenuItem)
	if notificationSoundEnabled {
		mSound = systray.AddMenuItemCheckbox("Sound", "Enable notification sound", settings.sound)
	}
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	menu := menuItems{mEnabled, mPause, mSound, mQuit}

	var ctx context.Context
	var ctxCancel context.CancelFunc

	for {
		select {
		case <-mEnabled.ClickedCh:
			if mEnabled.Checked() {
				mainCtxCancel()

				mEnabled.Uncheck()
				mPause.Disable()
			} else {
				runTwentyTwentyTwenty(notifier, &settings)

				mEnabled.Check()
				mPause.Enable()
			}
		case <-mPause.ClickedCh:
			if mPause.Checked() {
				mainCtxCancel() // make sure the current twenty-twenty-twenty goroutine stopped
				ctxCancel()     // cancel the current pause if it is running
				runTwentyTwentyTwenty(notifier, &settings)

				mEnabled.Enable()
				mPause.Uncheck()
			} else {
				ctx, ctxCancel = context.WithCancel(context.Background())
				go resumeTwentyTwentyTwentyAfter(ctx, ctxCancel, &settings, &menu)

				mEnabled.Disable()
				mPause.Check()
			}
		case <-mSound.ClickedCh:
			if mSound.Checked() {
				settings.sound = false

				mSound.Uncheck()
			} else {
				err := initNotification()
				if err != nil {
					log.Fatalf("Error while initialising sound: %v\n", err)
				}
				settings.sound = true

				mSound.Check()
			}
		case <-mQuit.ClickedCh:
			systray.Quit()
		}
	}
}

func onExit() {}
