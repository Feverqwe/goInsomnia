package internal

import (
	"fmt"
	"goInsomnia/assets"
	"runtime"

	"github.com/getlantern/systray"
)

var icon []byte

func TrayIcon(pc *PowerControl) {
	if icon == nil {
		data, err := assets.Asset("icon.ico")
		if err != nil {
			panic(err)
		}
		icon = data
	}

	onRun := func() {
		systray.SetTemplateIcon(icon, icon)
		systray.SetTooltip("GoInsomnia")

		reqTypeItem := map[uintptr]*systray.MenuItem{}

		var mLockExecuting *systray.MenuItem = nil
		if runtime.GOOS == "windows" {
			mLockExecuting = systray.AddMenuItemCheckbox("Executing", "Executing", false)
			reqTypeItem[EXECUTING] = mLockExecuting
		}

		mLockDisplay := systray.AddMenuItemCheckbox("Display", "Display", false)
		reqTypeItem[DISPLAY] = mLockDisplay

		mLockSystem := systray.AddMenuItemCheckbox("System", "System", false)
		reqTypeItem[SYSTEM] = mLockSystem

		var mLockAwayMode *systray.MenuItem = nil
		if runtime.GOOS == "windows" {
			mLockAwayMode = systray.AddMenuItemCheckbox("AwayMode", "AwayMode", false)
			reqTypeItem[AWAYMODE] = mLockAwayMode
		}

		mQuit := systray.AddMenuItem("Quit", "Quit")

		syncMenu := func() {
			for reqType, item := range reqTypeItem {
				enabled := pc.State[reqType]
				if enabled != item.Checked() {
					if enabled {
						item.Check()
					} else {
						item.Uncheck()
					}
				}
			}
		}

		onClick := func(item *systray.MenuItem, cb func(enabled bool) error) {
			enabled := !item.Checked()
			if err := cb(enabled); err == nil {
				syncMenu()
			} else {
				fmt.Println("Change state error", err)
			}
		}

		go func() {
			if runtime.GOOS == "windows" {
				onClick(mLockExecuting, pc.Executing)
			} else {
				onClick(mLockSystem, pc.System)
			}
		}()

		go func() {
			for {
				select {
				case <-mQuit.ClickedCh:
					systray.Quit()
				case <-mLockDisplay.ClickedCh:
					onClick(mLockDisplay, pc.Display)
				case <-mLockSystem.ClickedCh:
					onClick(mLockSystem, pc.System)
				}
			}
		}()

		if runtime.GOOS == "windows" {
			go func() {
				for {
					select {
					case <-mLockExecuting.ClickedCh:
						onClick(mLockExecuting, pc.Executing)
					case <-mLockAwayMode.ClickedCh:
						onClick(mLockAwayMode, pc.AwayMode)
					}
				}
			}()
		}
	}

	onExit := func() {

	}

	systray.Run(onRun, onExit)
}
