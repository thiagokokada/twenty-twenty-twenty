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

//go:embed assets/eye_light.png
var data []byte

func onReady() {
	systray.SetIcon(data)
	systray.SetTooltip("TwentyTwentyTwenty")
	mEnabled := systray.AddMenuItemCheckbox("Enabled", "Enable twenty-twenty-twenty", true)
	mPause := systray.AddMenuItemCheckbox("Pause for 1 hour", "Pause twenty-twenty-twenty for 1 hour", false)
	mSound := new(systray.MenuItem)
	if notificationSoundEnabled {
		mSound = systray.AddMenuItemCheckbox("Sound", "Enable notification sound", *notificationSound)
	}
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
	var ctx context.Context
	var cancel context.CancelFunc

	for {
		select {
		case <-mEnabled.ClickedCh:
			if mEnabled.Checked() {
				mainCtxCancel()
				mEnabled.Uncheck()
				mPause.Disable()
			} else {
				runTwentyTwentyTwenty()
				mEnabled.Check()
				mPause.Enable()
			}
		case <-mPause.ClickedCh:
			if mPause.Checked() {
				mainCtxCancel() // make sure the loop stopped
				cancel()
				runTwentyTwentyTwenty()
				mEnabled.Enable()
				mPause.Uncheck()
			} else {
				log.Println("Pausing twenty-twenty-twenty for 1 hour...")
				ctx, cancel = context.WithCancel(context.Background())
				go func(ctx context.Context) {
					mainCtxCancel()
					timer := time.NewTimer(time.Hour)
					select {
					case <-timer.C:
						notification := sendNotification(
							notifier,
							"Resuming 20-20-20",
							fmt.Sprintf("You will see a notification every %.f minutes(s)", frequency.Minutes()),
							notificationSound,
						)
						cancelNotificationAfter(notification, duration, notificationSound)
						runTwentyTwentyTwenty()
						mEnabled.Enable()
						mPause.Uncheck()
					case <-ctx.Done():
					}
					cancel()
				}(ctx)
				mEnabled.Disable()
				mPause.Check()
			}
		case <-mSound.ClickedCh:
			if mSound.Checked() {
				*notificationSound = false
				mSound.Uncheck()
			} else {
				err := initNotification()
				if err != nil {
					log.Fatalf("Error while initialising sound: %v\n", err)
				}
				*notificationSound = true
				mSound.Check()
			}
		case <-mQuit.ClickedCh:
			systray.Quit()
		}
	}
}

func onExit() {}
