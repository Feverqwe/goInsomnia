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
	IsDisplay   bool
	IsExecuting bool
	IsSystem    bool
	IsAwayMode  bool
	ctx         uintptr
}

func (self *PowerControl) Executing(enabled bool) error {
	var err error
	if enabled != self.IsExecuting {
		if enabled {
			err = self.set(EXECUTING)
		} else {
			err = self.clear(EXECUTING)
		}
	}
	if err == nil {
		self.IsExecuting = enabled
	}
	return err
}

func (self *PowerControl) Display(enabled bool) error {
	var err error
	if enabled != self.IsDisplay {
		if enabled {
			err = self.set(DISPLAY)
		} else {
			err = self.clear(DISPLAY)
		}
	}
	if err == nil {
		self.IsDisplay = enabled
	}
	return err
}

func (self *PowerControl) System(enabled bool) error {
	var err error
	if enabled != self.IsSystem {
		if enabled {
			err = self.set(SYSTEM)
		} else {
			err = self.clear(SYSTEM)
		}
	}
	if err == nil {
		self.IsSystem = enabled
	}
	return err
}

func (self *PowerControl) Awaymode(enabled bool) error {
	var err error
	if enabled != self.IsAwayMode {
		if enabled {
			err = self.set(AWAYMODE)
		} else {
			err = self.clear(AWAYMODE)
		}
	}
	if err == nil {
		self.IsAwayMode = enabled
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

type ULONG uint32
type DWORD uint32

type SimpleReasonString struct {
	SimpleReasonString *uint16
}

type REASON_CONTEXT struct {
	Version ULONG
	Flags   DWORD
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

	return powerControl
}
