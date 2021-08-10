package internal

import (
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

const SYSTEM = uintptr(0x1)
const EXECUTING = uintptr(0x3)
const AWAYMODE = uintptr(0x2)
const DISPLAY = uintptr(0x1 & 0x4)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var powerCreateRequest = kernel32.NewProc("PowerCreateRequest")
var powerSetRequest = kernel32.NewProc("PowerSetRequest")
var powerClearRequest = kernel32.NewProc("PowerClearRequest")

type PowerControl struct {
	State map[uintptr]bool
	ctx   uintptr
}

func (self *PowerControl) Executing(enabled bool) error {
	return self.change(EXECUTING, enabled)
}

func (self *PowerControl) Display(enabled bool) error {
	return self.change(DISPLAY, enabled)
}

func (self *PowerControl) System(enabled bool) error {
	return self.change(SYSTEM, enabled)
}

func (self *PowerControl) AwayMode(enabled bool) error {
	return self.change(AWAYMODE, enabled)
}

func (self *PowerControl) change(reqType uintptr, enabled bool) error {
	var err error
	if enabled != self.State[reqType] {
		if enabled {
			err = self.set(reqType)
		} else {
			err = self.clear(reqType)
		}
	}
	if err == nil {
		self.State[reqType] = enabled
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

	powerControl := &PowerControl{}
	powerControl.ctx = pCtx
	powerControl.State = map[uintptr]bool{}

	return powerControl
}
