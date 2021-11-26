package internal

import (
	"fmt"
	"goInsomnia/assets"
	"reflect"
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

		var mLockArr []*systray.MenuItem
		var mLockChannels []reflect.SelectCase

		syncMenu := func() {
			for index, powerType := range pc.types {
				menuItem := mLockArr[index]
				enabled := powerType.state
				if enabled != menuItem.Checked() {
					if enabled {
						menuItem.Check()
					} else {
						menuItem.Uncheck()
					}
				}
			}
		}

		onClick := func(index int) {
			menuItem := mLockArr[index]
			powerType := pc.types[index]
			state := !menuItem.Checked()
			if err := pc.setState(powerType.id, state); err == nil {
				syncMenu()
			} else {
				fmt.Println("Change state error", err)
			}
		}

		for _, powerType := range pc.types {
			mLock := systray.AddMenuItemCheckbox(powerType.title, powerType.tooltip, false)
			mLockArr = append(mLockArr, mLock)
			selectCase := reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(mLock.ClickedCh)}
			mLockChannels = append(mLockChannels, selectCase)
		}

		mQuit := systray.AddMenuItem("Quit", "Quit")

		go func() {
			if runtime.GOOS == "windows" {
				onClick(0)
			} else {
				onClick(1)
			}
		}()

		go func() {
			for {
				select {
				case <-mQuit.ClickedCh:
					systray.Quit()
				}
			}
		}()

		go func() {
			for {
				index, _, _ := reflect.Select(mLockChannels)
				onClick(index)
			}
		}()
	}

	onExit := func() {

	}

	systray.Run(onRun, onExit)
}
