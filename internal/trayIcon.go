package internal

import (
	"fmt"
	"goInsomnia/assets"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/getlantern/systray"
)

var icon []byte
var disabledIcon []byte

var minutesPreset = []int{5, 10, 15, 30, 60, 120, 240, 360, 480}

func TrayIcon(pc *PowerControl) {
	if icon == nil {
		data, err := assets.Asset("icon.ico")
		if err != nil {
			panic(err)
		}
		icon = data
	}

	if disabledIcon == nil {
		data, err := assets.Asset("disabled.ico")
		if err != nil {
			panic(err)
		}
		disabledIcon = data
	}

	onRun := func() {
		systray.SetTemplateIcon(disabledIcon, disabledIcon)
		systray.SetTooltip("GoInsomnia")

		var mLockArr []*systray.MenuItem
		var mLockChannels []reflect.SelectCase

		var mMinutesChannels []reflect.SelectCase

		syncMenu := func() {
			hasEnabled := false
			for index, powerType := range pc.types {
				menuItem := mLockArr[index]
				enabled := powerType.state
				if enabled {
					hasEnabled = true
				}
				if enabled != menuItem.Checked() {
					if enabled {
						menuItem.Check()
					} else {
						menuItem.Uncheck()
					}
				}
			}
			if hasEnabled {
				systray.SetTemplateIcon(icon, icon)
			} else {
				systray.SetTemplateIcon(disabledIcon, disabledIcon)
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

		subMinutes := systray.AddMenuItem("Turn off after...", "Turn off after...")

		var timer *time.Timer
		stopTimer := func() {
			if timer != nil && timer.Stop() {
				timerItem.Hide()
				subMinutes.Show()
			}
		}
		onTimer := func() {
			for _, powerType := range pc.types {
				if powerType.state {
					pc.setState(powerType.id, false)
				}
			}
			syncMenu()
		}
		setTimer := func(minutes int) {
			stopTimer()
			duration := time.Duration(minutes) * time.Minute
			ct := time.Now()
			ct = ct.Add(duration)
			timer = time.AfterFunc(duration, onTimer)
			format := "15:04"
			if minutes > 24*60 {
				format = "Jan 2 15:04"
			}
			timerItem.SetTitle("Until " + ct.Format(format))
			subMinutes.Hide()
			timerItem.Show()
		}

		onMinutesClick := func(index int) {
			minutes := minutesPreset[index]
			setTimer(minutes)
		}

		for _, minutes := range minutesPreset {
			title := formatMinutes(minutes)
			mMinutes := subMinutes.AddSubMenuItem(title, title)
			selectCase := reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(mMinutes.ClickedCh)}
			mMinutesChannels = append(mMinutesChannels, selectCase)
		}

		mSetMinutes := subMinutes.AddSubMenuItem("Set minutes", "Set minues")

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
				case <-mSetMinutes.ClickedCh:
					minutesStr, err := ShowEntry("Set minutes", "Enter the number of minutes:", "60")
					var minutes int
					if err == nil {
						minutes, err = strconv.Atoi(minutesStr)
					}
					if err != nil {
						fmt.Println("Set minutes error", err)
					} else {
						setTimer(minutes)
					}
				case <-timerItem.ClickedCh:
					stopTimer()
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
