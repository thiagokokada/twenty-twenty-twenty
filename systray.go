//go:build !nosystray
// +build !nosystray

package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"fyne.io/systray"

	"github.com/thiagokokada/twenty-twenty-twenty/core"
	"github.com/thiagokokada/twenty-twenty-twenty/sound"
)

const systrayEnabled bool = true

func withMutex(mu *sync.Mutex, f func()) {
	mu.Lock()
	defer mu.Unlock()
	f()
}

func onReady() {
	systray.SetIcon(systrayIcon)
	systray.SetTooltip("TwentyTwentyTwenty")
	mEnabled := systray.AddMenuItemCheckbox("Enabled", "Enable twenty-twenty-twenty", true)
	mPause := systray.AddMenuItemCheckbox(
		fmt.Sprintf("Pause for %.f hour", settings.Pause.Hours()),
		fmt.Sprintf("Pause twenty-twenty-twenty for %.f hour", settings.Pause.Hours()),
		false,
	)
	mSound := new(systray.MenuItem)
	if sound.Enabled {
		mSound = systray.AddMenuItemCheckbox("Sound", "Enable notification sound", settings.Sound)
	}
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quit the whole app")

	var pauseCtx context.Context
	var cancelPauseCtx context.CancelFunc
	var mu sync.Mutex

	for {
		select {
		case <-mEnabled.ClickedCh:
			if mEnabled.Checked() {
				core.Stop()

				withMutex(&mu, func() {
					mEnabled.Uncheck()
					mPause.Disable()
				})
			} else {
				core.Start(&settings, optional)

				withMutex(&mu, func() {
					mEnabled.Check()
					mPause.Enable()
				})
			}
		case <-mPause.ClickedCh:
			if mPause.Checked() {
				cancelPauseCtx() // cancel the current pause if it is running
				core.Start(&settings, optional)

				withMutex(&mu, func() {
					mEnabled.Enable()
					mPause.Uncheck()
				})
			} else {
				pauseCtx, cancelPauseCtx = context.WithCancel(context.Background())
				go func() {
					defer cancelPauseCtx()
					core.Pause(
						pauseCtx, &settings, optional,
						func() { withMutex(&mu, func() { mPause.Disable() }) }, // blocking pause button to avoid concurrency issue
						func() {
							withMutex(&mu, func() {
								mEnabled.Enable()
								mPause.Uncheck()
								mPause.Enable()
							})
						},
					)
				}()

				withMutex(&mu, func() {
					mEnabled.Disable()
					mPause.Check()
				})
			}
		case <-mSound.ClickedCh:
			if mSound.Checked() {
				settings.Sound = false

				withMutex(&mu, func() { mSound.Uncheck() })
			} else {
				err := sound.Init()
				if err != nil {
					log.Fatalf("Error while initialising sound: %v\n", err)
				}
				settings.Sound = true

				withMutex(&mu, func() { mSound.Check() })
			}
		case <-mQuit.ClickedCh:
			systray.Quit()
		}
	}
}

func onExit() {}
