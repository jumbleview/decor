package screen

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/image/bmp"
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

// SetWallpaper decodes jpeg file into bmp and sets it as windows background
func SetWallpaper(jpegName, currentPath string) error {
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("SystemParametersInfoW")
	metrics := user32.NewProc("GetSystemMetrics")
	// Read Monitor X, Y resolution
	xNumber, _, _ := metrics.Call(SM_CXSCREEN)
	yNumber, _, _ := metrics.Call(SM_CYSCREEN)

	monitorWidth := int(*(*int)(unsafe.Pointer(&xNumber)))
	monitorHeight := int(*(*int)(unsafe.Pointer(&yNumber)))

	fconfig, err := os.Open(jpegName)

	if err != nil {
		return err
	}
	// Read picture format and pixel size
	config, format, err := image.DecodeConfig(fconfig)
	fconfig.Close()
	if err != nil {
		return err
	}

	pictureShape := float64(config.Width) / float64(config.Height)
	monitorShape := float64(monitorWidth) / float64(monitorHeight)
	cropWigth := config.Width
	cropHeight := monitorHeight * cropWigth / monitorWidth
	cropX := 0
	cropY := (config.Height - cropHeight) / 2
	if pictureShape > monitorShape {
		cropHeight = config.Height
		cropWigth = monitorWidth * cropHeight / monitorHeight
		cropX = (config.Width - cropWigth) / 2
		cropY = 0
	}
	fmt.Printf("Image  type: %s; size: %d*%d\n", format, config.Width, config.Height)
	fmt.Printf("Monitor size: %d*%d\n", monitorWidth, monitorHeight)
	fmt.Printf("Crop: %d, %d; %d %d\n", cropWigth, cropHeight, cropX, cropY)

	f, err := os.Open(jpegName)
	if err != nil {
		return err
	}
	defer f.Close()
	if format != "jpeg" {
		return fmt.Errorf("image format is %s != jpeg", format)
	}
	// Decode jpeg file into image
	jpegr := bufio.NewReader(f)
	imgFull, err := jpeg.Decode(jpegr)
	if err != nil {
		return err
	}

	// Crop image to fit monitor screen to avoid image distortion
	var imgCrop image.Image
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}
	simg, ok := imgFull.(subImager)
	rect := image.Rect(cropX, cropY, cropWigth, cropHeight)

	if ok {
		imgCrop = simg.SubImage(rect)
	}

	bmpName := filepath.Join(currentPath, "jumbleview_current.bmp")
	fbmp, err := os.Create(bmpName)
	if err != nil {
		return err
	}
	defer fbmp.Close()
	// encode image as bmp file
	bmpw := bufio.NewWriter(fbmp)
	if ok {
		err = bmp.Encode(bmpw, imgCrop)
	} else {
		err = bmp.Encode(bmpw, imgFull)
	}
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", bmpName)
	address, err := syscall.UTF16PtrFromString(bmpName)
	if err == nil {
		proc.Call(
			spiSETDESKWALLPAPER,
			0,
			uintptr(unsafe.Pointer(address)),
			spifUPDATEINIFILE,
		)
	}
	return err
}
