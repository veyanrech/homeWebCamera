package image

import (
	"image"
	"os"

	_ "image/jpeg"
	_ "image/png"
)

const darkPixelPercentage = 80
const darkPixelThreshold = 38

// IsImageBlack checks if the image is black
// i dont close the file, because it should be closed in the calling function
func IsImageBlack(src *os.File) (r bool, err error) {

	img, _, err := image.Decode(src)
	// fmt.Println(frmt)

	if err != nil {
		return false, err
	}

	darkPixelCount := 0
	totalPixels := 0

	bounds := img.Bounds()

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			totalPixels++
			r, g, b, _ := img.At(x, y).RGBA()
			r, g, b = r>>8, g>>8, b>>8
			if int(r) < darkPixelThreshold && int(g) < darkPixelThreshold && int(b) < darkPixelThreshold {
				darkPixelCount++
			}
		}
	}

	if darkPixelCount*100/totalPixels >= darkPixelPercentage {
		return true, nil
	}

	return false, nil
}
