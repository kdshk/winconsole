package winconsole

import (
	"runtime"
	"syscall"
	"unsafe"
	"errors"
)

const (
	// STD INPUT Flags
	ENABLE_ECHO_INPUT 	   = 0x0004
	ENABLE_INSERT_MODE     = 0x0020
	ENABLE_LINE_INPUT	   = 0x0002
	ENABLE_MOUSE_INPUT     = 0x0010
	ENABLE_PROCESSED_INPUT = 0x0001
	ENABLE_QUICK_EDIT_MODE = 0x0040
	ENABLE_WINDOW_INPUT    = 0x0008
	// Required to enable or disable extended flags	
	ENABLE_EXTENDED_FLAGS  = 0x0080	
	ENABLE_VIRTUAL_TERMINAL_INPUT = 0x0200
	// STD OUTPUT Flags
	ENABLE_PROCESSED_OUTPUT = 0x0001
	ENABLE_WRAP_AT_EOL_OUTPUT = 0x0002
	ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	DISABLE_NEWLINE_AUTO_RETURN = 0x0008
	ENABLE_LVB_GRID_WORLDWIDE = 0x0010
	
	STD_INPUT_HANDLE	   = -10
	STD_OUTPUT_HANDLE	   = -11
	STD_ERROR_HANDLE	   = -12
	//
	STD_INPUT_NO_QUICK_MODE = 0
)

/*
  Get Windows Console Flags for std input,output and error
*/
func GetConsoleFlag(handle int) (error, uint32) {
	lpMode := uint32(0)
	sysType := runtime.GOOS
	if sysType == "windows" {	
		stdHandle, err := syscall.GetStdHandle(handle)
		if err != nil{
			return err, 0
		}
		kernel32 , err := syscall.LoadLibrary("kernel32.dll")
		if err != nil{
			return err, 0
		}
		// release the libary
		defer syscall.FreeLibrary(kernel32)
		procGetConsoleMode , err := syscall.GetProcAddress(kernel32, "GetConsoleMode")
		if err != nil{
			return err, 0
		}
		r , _ , err := syscall.Syscall(uintptr(procGetConsoleMode), 2,
			uintptr(stdHandle), 
			uintptr(unsafe.Pointer(&lpMode)),
			0)
		if err != nil{
			return err, 0
		}
		if r == 0 {
			return errors.New("call get console mode failed"), lpMode
		} 
		return nil, lpMode
	} else {
		return errors.New("not a windows os"), 0
	}
}

/*
  Set Windows Console Flags for std input,output and error
*/
func SetConsoleFlag(flag uint32, handle int) error {
	sysType := runtime.GOOS
	if sysType == "windows" {	
		stdHandle, err := syscall.GetStdHandle(handle)
		if err != nil{
			return err
		}
		kernel32 , err := syscall.LoadLibrary("kernel32.dll")
		if err != nil{
			return err
		}
		// release the libary
		defer syscall.FreeLibrary(kernel32)
		procSetConsoleMode , err := syscall.GetProcAddress(kernel32, "SetConsoleMode")
		if err != nil{
			return err
		}
		r , _ , err := syscall.Syscall(uintptr(procSetConsoleMode), 2,
			uintptr(stdHandle), 
			uintptr(flag),
			0)
		if err != nil{
			return err
		}
		if r == 0 {
			return errors.New("call get console mode failed")
		} 
		return nil
	} else {
		return errors.New("not a windows os")
	}
}
/*
	initialize the windows console mode to disable quick edit mode 
	to avoid console stuck by enter quick edit mode on left click
 */
func DisableConsoleQuickEditMode() {
	sysType := runtime.GOOS
	if sysType == "windows" {
		// get std input handle to set options 
		// STD_INPUT_HANDLE (DWORD) -10
		stdHandle, _ := syscall.GetStdHandle(STD_INPUT_HANDLE)
		setConsoleQuickEditMode(stdHandle, false)
	}
}

func setConsoleQuickEditMode(handle syscall.Handle, enable bool) error {
	// load kernel libary
	lpMode := uint32(0)
	kernel32 , err := syscall.LoadLibrary("kernel32.dll")
	if err != nil{
		return err
	}
	// release the libary
	defer syscall.FreeLibrary(kernel32)

	procGetConsoleMode , err := syscall.GetProcAddress(kernel32, "GetConsoleMode")
	if err != nil{
		return err
	}
	procSetConsoleMode , err := syscall.GetProcAddress(kernel32, "SetConsoleMode")
	if err != nil{
		return err
	}
	r , _ , err := syscall.Syscall(uintptr(procGetConsoleMode), 2,
			uintptr(handle), 
			uintptr(unsafe.Pointer(&lpMode)),
			0)
	if r != 0 {
		lpMode |= ENABLE_EXTENDED_FLAGS	
		if enable {
			lpMode |= ENABLE_QUICK_EDIT_MODE
		} else {
			lpMode &= ^uint32(ENABLE_QUICK_EDIT_MODE)
		}
			r , _ , err = syscall.Syscall(uintptr(procSetConsoleMode), 2,
			uintptr(handle), 
			uintptr(lpMode),
			0)
	}
	return err
}


