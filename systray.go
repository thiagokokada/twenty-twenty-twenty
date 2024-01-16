package main

import (
	"log"

	"fyne.io/systray"
)

func onReady() {
	systray.SetIcon(data)
	systray.SetTooltip("TwentyTwentyTwenty")
	mSound := systray.AddMenuItemCheckbox("Sound", "Enable notification sound", *notificationSound)
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	for {
		select {
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
