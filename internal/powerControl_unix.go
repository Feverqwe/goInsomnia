package internal

import (
	"log"

	"github.com/caseymrm/go-caffeinate"
)

const SYSTEM = uintptr(0x1)
const EXECUTING = uintptr(0x3)
const AWAYMODE = uintptr(0x2)
const DISPLAY = uintptr(0x1 & 0x4)

type PowerControl struct {
	State map[uintptr]bool
	ch    chan int
}

func (self *PowerControl) Executing(enabled bool) error {
	return self.change(EXECUTING, false)
}

func (self *PowerControl) Display(enabled bool) error {
	return self.change(DISPLAY, enabled)
}

func (self *PowerControl) System(enabled bool) error {
	return self.change(SYSTEM, enabled)
}

func (self *PowerControl) AwayMode(enabled bool) error {
	return self.change(AWAYMODE, false)
}

func (self *PowerControl) change(reqType uintptr, enabled bool) error {
	self.State[reqType] = enabled
	self.ch <- 1
	return nil
}

func GetPowerControl() *PowerControl {
	powerControl := &PowerControl{
		ch:    make(chan int),
		State: map[uintptr]bool{},
	}
	var c *caffeinate.Caffeinate
	go func() {
		for {
			<-powerControl.ch
			log.Println("caffeinate: changes")
			someone := false
			system := false
			display := false
			for key, value := range powerControl.State {
				if value {
					someone = true
				}
				if key == SYSTEM {
					system = value
				}
				if key == DISPLAY {
					display = value
				}
			}

			if c != nil && c.Running() {
				log.Println("caffeinate: stop")
				err := c.Stop()
				if err != nil {
					log.Println("caffeinate: stop error", err)
				}
			}
			if someone {
				c = &caffeinate.Caffeinate{
					System:  system,
					Display: display,
				}
				log.Println("caffeinate: start")
				c.Start()
			}
		}
	}()
	return powerControl
}
