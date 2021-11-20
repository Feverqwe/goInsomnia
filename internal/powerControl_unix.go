//go:build darwin

package internal

import (
	"log"

	"github.com/caseymrm/go-caffeinate"
)

type PowerType struct {
	id      int
	title   string
	tooltip string
	state   bool
}

type PowerControl struct {
	ch     chan int
	types  []*PowerType
	idType map[int]*PowerType
}

func (self *PowerControl) setState(id int, state bool) error {
	powerType := self.idType[id]
	powerType.state = state
	self.ch <- 1
	return nil
}

func GetPowerControl() *PowerControl {
	var types []*PowerType
	types = append(types, &PowerType{
		id:      1,
		title:   "Display",
		tooltip: "Prevent the display from sleeping.",
	})
	types = append(types, &PowerType{
		id:      2,
		title:   "Idle System",
		tooltip: "Prevent the system from idle sleeping.",
	})
	types = append(types, &PowerType{
		id:      3,
		title:   "Idle Disk",
		tooltip: "Prevent the disk from idle sleeping.",
	})
	types = append(types, &PowerType{
		id:      4,
		title:   "System",
		tooltip: "Prevent the system from sleeping. Valid only on AC power.",
	})

	idType := make(map[int]*PowerType)
	for _, powerType := range types {
		idType[powerType.id] = powerType
	}

	powerControl := &PowerControl{
		ch:     make(chan int),
		types:  types,
		idType: idType,
	}

	var c *caffeinate.Caffeinate
	go func() {
		for {
			<-powerControl.ch
			log.Println("caffeinate: changes")
			someone := false
			display := false
			idleSystem := false
			idleDisk := false
			system := false
			for _, powerType := range powerControl.types {
				value := powerType.state
				if value {
					someone = true
				}
				switch powerType.id {
				case 1:
					display = value
				case 2:
					idleSystem = value
				case 3:
					idleDisk = value
				case 4:
					system = value
				}
			}

			if c != nil && c.Running() {
				log.Println("caffeinate: stop")
				err := c.Stop()
				if err != nil && err.Error() != "signal: killed" {
					log.Println("caffeinate: stop error", err.Error())
				}
			}
			if someone {
				c = &caffeinate.Caffeinate{
					Display:    display,
					IdleSystem: idleSystem,
					IdleDisk:   idleDisk,
					System:     system,
				}
				log.Println("caffeinate: start")
				c.Start()
			}
		}
	}()
	return powerControl
}
