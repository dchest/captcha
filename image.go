package captcha

import (
	"image"
	"image/png"
	"io"
	"os"
	"rand"
)


func NewImage(numbers []byte) *image.NRGBA {
	w := numberWidth * (dotSize + 3) * len(numbers)
	h := numberHeight * (dotSize + 5)
	img := image.NewNRGBA(w, h)
	color := image.NRGBAColor{uint8(rand.Intn(50)), uint8(rand.Intn(50)), uint8(rand.Intn(128)), 0xFF}
	fillWithCircles(img, color, 40, 4)
	x := rand.Intn(dotSize)
	y := 0
	setRandomBrightness(&color, 180)
	for _, n := range numbers {
		y = rand.Intn(dotSize * 4)
		drawNumber(img, font[n], x, y, color)
		x += dotSize*numberWidth + rand.Intn(maxSkew) + 8
	}
	drawCirclesLine(img, color)
	return img
}

func EncodeNewImage(w io.Writer) (numbers []byte, err os.Error) {
	numbers = randomNumbers()
	err = png.Encode(w, NewImage(numbers))
	return
}

