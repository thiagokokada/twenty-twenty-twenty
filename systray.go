//go:build !nosystray
// +build !nosystray

package main

import (
	_ "embed"
	"log"

	"fyne.io/systray"
)

//go:embed assets/eye_light.png
var data []byte

func onReady() {
	systray.SetIcon(data)
	systray.SetTooltip("TwentyTwentyTwenty")
	mEnabled := systray.AddMenuItemCheckbox("Enabled", "Enable twenty-twenty-twenty", true)
	mSound := new(systray.MenuItem)
	if notificationSoundEnabled {
		mSound = systray.AddMenuItemCheckbox("Sound", "Enable notification sound", *notificationSound)
	}
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	for {
		select {
		case <-mEnabled.ClickedCh:
			if mEnabled.Checked() {
				ctxCancel()
				mEnabled.Uncheck()
			} else {
				runTwentyTwentyTwenty()
				mEnabled.Check()
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
