//go:build !nosystray
// +build !nosystray

package main

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"fyne.io/systray"

	"github.com/thiagokokada/twenty-twenty-twenty/internal/core"
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
	if optional.Sound {
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
				slog.DebugContext(core.Ctx(), "Enable button unchecked")
				core.Stop()

				withMutex(&mu, func() {
					mEnabled.Uncheck()
					mPause.Disable()
				})
			} else {
				slog.DebugContext(core.Ctx(), "Enable button checked")
				core.Start(&settings, optional)

				withMutex(&mu, func() {
					mEnabled.Check()
					mPause.Enable()
				})
			}
		case <-mPause.ClickedCh:
			if pauseCtx != nil {
				slog.DebugContext(core.Ctx(), "Cancelling current pause")
				cancelPauseCtx()
			}
			if mPause.Checked() {
				slog.DebugContext(core.Ctx(), "Pause button unchecked")
				core.Start(&settings, optional)

				withMutex(&mu, func() {
					mEnabled.Enable()
					mPause.Uncheck()
				})
			} else {
				slog.DebugContext(core.Ctx(), "Pause button checked")
				pauseCtx, cancelPauseCtx = context.WithCancel(context.Background())
				go func() {
					defer cancelPauseCtx()
					core.Pause(
						pauseCtx, &settings, optional,
						func() {
							withMutex(&mu, func() {
								slog.DebugContext(core.Ctx(), "Calling pause callback")
								mEnabled.Enable()
								mPause.Uncheck()
							})
						},
						nil,
					)
				}()

				withMutex(&mu, func() {
					mEnabled.Disable()
					mPause.Check()
				})
			}
		case <-mSound.ClickedCh:
			if mSound.Checked() {
				slog.DebugContext(core.Ctx(), "Sound button unchecked")
				settings.Sound = false

				withMutex(&mu, func() { mSound.Uncheck() })
			} else {
				slog.DebugContext(core.Ctx(), "Sound button checked")
				settings.Sound = true

				withMutex(&mu, func() { mSound.Check() })
			}
		case <-mQuit.ClickedCh:
			slog.DebugContext(core.Ctx(), "Quit button clicked")
			systray.Quit()
		}
	}
}

func onExit() {
	core.Stop()
}
