package main

import (
	"fmt"
	"goInsomnia/internal"
)

func main() {
	if _, err := internal.CreateMutex("GoInsomnia"); err != nil {
		panic(err)
	}

	state := internal.State{
		Executing: true,
	}

	ch := make(chan string)

	internal.TrayIcon(&state, ch)

	pc := internal.GetPowerControl()

	go func() { ch <- "sync" }()

	for {
		v := <-ch
		switch v {
		case "sync":
			err := pc.Executing(state.Executing)
			if err != nil {
				fmt.Println("Set state error", err)
			}
			err = pc.Display(state.Display)
			if err != nil {
				fmt.Println("Set state error", err)
			}
			err = pc.System(state.System)
			if err != nil {
				fmt.Println("Set state error", err)
			}
			err = pc.Awaymode(state.AwayMode)
			if err != nil {
				fmt.Println("Set state error", err)
			}
		}
	}
}
