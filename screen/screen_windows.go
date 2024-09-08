package screen

import (
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

// MakeConsole allocates console for the app built with -H=windowsgui flag
func MakeConsole() (int, error) {
	modkernel32 := syscall.NewLazyDLL("kernel32.dll")
	procAllocConsole := modkernel32.NewProc("AllocConsole")
	r0, _, err0 := syscall.SyscallN(procAllocConsole.Addr())
	if r0 == 0 { // Allocation failed, probably process already has a console
		return 1, err0
	}
	hout, err1 := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	hin, err2 := syscall.GetStdHandle(syscall.STD_INPUT_HANDLE)
	if err1 != nil || err2 != nil { // nowhere to print the error
		return 2, err0
	}
	os.Stdout = os.NewFile(uintptr(hout), "/dev/stdout")
	os.Stdin = os.NewFile(uintptr(hin), "/dev/stdin")
	return 0, err0
}

const (
	spiSETDESKWALLPAPER = 0x14
	spifUPDATEINIFILE   = 0x2
	SM_CXSCREEN         = 0
	SM_CYSCREEN         = 1
)

// SetWallpaper decodes jpeg file into bmp and after cropping sets it as windows background
func SetWallpaper(jpegName, currentPath string) error {
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("SystemParametersInfoW")
	metrics := user32.NewProc("GetSystemMetrics")

	// Read Monitor X, Y resolution
	xNumber, _, _ := metrics.Call(SM_CXSCREEN)
	yNumber, _, _ := metrics.Call(SM_CYSCREEN)

	monitorWidth := int(*(*int)(unsafe.Pointer(&xNumber)))
	monitorHeight := int(*(*int)(unsafe.Pointer(&yNumber)))

	bmpName := filepath.Join(currentPath, "jumbleview_current.bmp")
	// Convert JPG file into image, crop it, and save as BMP file
	err := EncodeCropDecode(jpegName, bmpName, monitorWidth, monitorHeight)
	if err != nil {
		return err
	}

	// Set BMP file as background
	address, err := syscall.UTF16PtrFromString(bmpName)
	if err != nil {
		return err
	}
	_, _, err = proc.Call(
		spiSETDESKWALLPAPER,
		0,
		uintptr(unsafe.Pointer(address)),
		spifUPDATEINIFILE,
	)
	return err
}
