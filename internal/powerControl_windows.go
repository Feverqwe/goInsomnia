//go:build windows

package internal

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var powerCreateRequest = kernel32.NewProc("PowerCreateRequest")
var powerSetRequest = kernel32.NewProc("PowerSetRequest")
var powerClearRequest = kernel32.NewProc("PowerClearRequest")

type PowerType struct {
	id      int
	title   string
	tooltip string
	state   bool
	reqType uintptr
}

type PowerControl struct {
	ctx    uintptr
	types  []*PowerType
	idType map[int]*PowerType
}

func (self *PowerControl) setState(id int, enabled bool) error {
	powerType := self.idType[id]
	var err error
	if enabled != powerType.state {
		reqType := powerType.reqType
		if enabled {
			err = self.set(reqType)
		} else {
			err = self.clear(reqType)
		}
	}
	if err == nil {
		powerType.state = enabled
	}
	return err
}

func (self *PowerControl) set(reqType uintptr) error {
	_, errNum, status := powerSetRequest.Call(self.ctx, reqType)
	if errNum != 0 {
		return status
	}
	return nil
}

func (self *PowerControl) clear(reqType uintptr) error {
	_, errNum, status := powerClearRequest.Call(self.ctx, reqType)
	if errNum != 0 {
		return status
	}
	return nil
}

type SimpleReasonString struct {
	SimpleReasonString *uint16
}

type REASON_CONTEXT struct {
	Version uint32
	Flags   uint32
	Reason  SimpleReasonString
}

func GetPowerControl() *PowerControl {
	var types []*PowerType
	types = append(types, &PowerType{
		id:      1,
		title:   "Executing",
		tooltip: "The calling process continues to run instead of being suspended or terminated by process lifetime management mechanisms. When and how long the process is allowed to run depends on the operating system and power policy settings.",
		reqType: uintptr(0x3),
	})
	types = append(types, &PowerType{
		id:      2,
		title:   "Display",
		tooltip: "The display remains on even if there is no user input for an extended period of time.",
		reqType: uintptr(0x1 & 0x4),
	})
	types = append(types, &PowerType{
		id:      3,
		title:   "System",
		tooltip: "The system continues to run instead of entering sleep after a period of user inactivity.",
		reqType: uintptr(0x1),
	})
	types = append(types, &PowerType{
		id:      4,
		title:   "Away mode",
		tooltip: "The system enters away mode instead of sleep in response to explicit action by the user. In away mode, the system continues to run but turns off audio and video to give the appearance of sleep. PowerRequestAwayModeRequired is only applicable on Traditional Sleep (S3) systems.",
		reqType: uintptr(0x2),
	})

	idType := make(map[int]*PowerType)
	for _, powerType := range types {
		idType[powerType.id] = powerType
	}

	sr, _ := windows.UTF16PtrFromString("Insomnia")
	reason := &REASON_CONTEXT{
		Version: 0,
		Flags:   0x1,
		Reason:  SimpleReasonString{sr},
	}
	pCtx, errNum, status := powerCreateRequest.Call(uintptr(unsafe.Pointer(reason)))
	if errNum != 0 {
		panic(status)
	}

	powerControl := &PowerControl{
		ctx:    pCtx,
		types:  types,
		idType: idType,
	}

	return powerControl
}
