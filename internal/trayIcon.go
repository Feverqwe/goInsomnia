package internal

import (
	"fmt"
	"goInsomnia/assets"
	"reflect"
	"runtime"
	"time"

	"github.com/getlantern/systray"
)

var icon []byte

var minutesPreset = []int{5, 10, 15, 30, 60, 120, 240, 360, 480}

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

		var mMinutesChannels []reflect.SelectCase

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

		timerItem := systray.AddMenuItem("", "")
		timerItem.Hide()

		subConfig := systray.AddMenuItem("Turn off after...", "Turn off after...")

		var timer *time.Timer
		stopTimer := func() {
			if timer != nil {
				timer.Stop()
			}
		}
		onTimer := func() {
			for id, powerType := range pc.idType {
				if powerType.state {
					pc.setState(id, false)
				}
			}
			syncMenu()
			timerItem.Hide()
			subConfig.Show()
		}
		setTimer := func(minutes int) {
			stopTimer()
			duration := time.Duration(minutes) * time.Minute
			ct := time.Now()
			ct = ct.Add(duration)
			timer = time.AfterFunc(duration, onTimer)
			timerItem.SetTitle("Until " + ct.Format("15:04"))
		}

		onTimerClick := func() {
			stopTimer()
			timerItem.Hide()
			subConfig.Show()
		}

		onMinutesClick := func(index int) {
			minutes := minutesPreset[index]
			setTimer(minutes)
			timerItem.Show()
			subConfig.Hide()
		}

		for _, minutes := range minutesPreset {
			title := formatMinutes(minutes)
			mMinutes := subConfig.AddSubMenuItem(title, title)
			selectCase := reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(mMinutes.ClickedCh)}
			mMinutesChannels = append(mMinutesChannels, selectCase)
		}

		systray.AddSeparator()

		for _, powerType := range pc.types {
			mLock := systray.AddMenuItemCheckbox(powerType.title, powerType.tooltip, false)
			mLockArr = append(mLockArr, mLock)
			selectCase := reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(mLock.ClickedCh)}
			mLockChannels = append(mLockChannels, selectCase)
		}

		systray.AddSeparator()

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
				case <-timerItem.ClickedCh:
					onTimerClick()
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

		go func() {
			for {
				index, _, _ := reflect.Select(mMinutesChannels)
				onMinutesClick(index)
			}
		}()
	}

	onExit := func() {

	}

	systray.Run(onRun, onExit)
}
