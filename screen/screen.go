package screen

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"os"

	"golang.org/x/image/bmp"
)

// EncodeCropDecode reads jpeg file, crops image to match monitor ratio, encode image as bmp file
func EncodeCropDecode(jpegName, bmpName string, monitorWidth, monitorHeight int) error {

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

	pictureRatio := float64(config.Width) / float64(config.Height)
	monitorRatio := float64(monitorWidth) / float64(monitorHeight)
	cropWigth := config.Width
	cropHeight := monitorHeight * cropWigth / monitorWidth
	cropX := 0
	cropY := (config.Height - cropHeight) / 2
	if pictureRatio > monitorRatio {
		cropHeight = config.Height
		cropWigth = monitorWidth * cropHeight / monitorHeight
		cropX = (config.Width - cropWigth) / 2
		cropY = 0
	}
	fmt.Printf("Image  type: %s\n", format)
	fmt.Printf("Image  size: %d*%d\n", config.Width, config.Height)
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
	return err
}
