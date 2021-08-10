package main

import (
	"goInsomnia/internal"
)

func main() {
	if _, err := internal.CreateMutex("GoInsomnia"); err != nil {
		panic(err)
	}

	pc := internal.GetPowerControl()

	internal.TrayIcon(pc)
}
