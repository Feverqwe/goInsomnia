package internal

import (
	"github.com/getlantern/systray"
	"goInsomnia/asserts"
	"os"
)

var icon []byte

func TrayIcon(state *State, callChan chan string) {
	if icon == nil {
		data, err := asserts.Asset("icon.ico")
		if err != nil {
			panic(err)
		}
		icon = data
	}

	onRun := func() {
		systray.SetTemplateIcon(icon, icon)
		systray.SetTitle("GoInsomnia")
		systray.SetTooltip("GoInsomnia")

		mLockExecuting := systray.AddMenuItemCheckbox("Executing", "Executing", state.Executing)
		mLockDisplay := systray.AddMenuItemCheckbox("Display", "Display", state.Display)
		mLockSystem := systray.AddMenuItemCheckbox("System", "System", state.System)
		mLockAwayMode := systray.AddMenuItemCheckbox("AwayMode", "AwayMode", state.AwayMode)

		mQuit := systray.AddMenuItem("Quit", "Quit")

		go func() {
			for {
				select {
				case <-mQuit.ClickedCh:
					systray.Quit()
					os.Exit(0)
				case <-mLockExecuting.ClickedCh:
					state.Executing = !state.Executing
					if state.Executing {
						mLockExecuting.Check()
					} else {
						mLockExecuting.Uncheck()
					}
					callChan <- "sync"
				case <-mLockDisplay.ClickedCh:
					state.Display = !state.Display
					if state.Display {
						mLockDisplay.Check()
					} else {
						mLockDisplay.Uncheck()
					}
					callChan <- "sync"
				case <-mLockSystem.ClickedCh:
					state.System = !state.System
					if state.System {
						mLockSystem.Check()
					} else {
						mLockSystem.Uncheck()
					}
					callChan <- "sync"
				case <-mLockAwayMode.ClickedCh:
					state.AwayMode = !state.AwayMode
					if state.AwayMode {
						mLockAwayMode.Check()
					} else {
						mLockAwayMode.Uncheck()
					}
					callChan <- "sync"
				}
			}
		}()
	}

	onExit := func() {}

	go func() {
		systray.Run(onRun, onExit)
	}()
}
