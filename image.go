package captcha

import (
	"image"
	"image/png"
	"io"
	"os"
	"rand"
)

const (
	maxSkew = 3
)

type CaptchaImage struct {
	*image.NRGBA
	primaryColor image.NRGBAColor
	numberWidth  int
	numberHeight int
	dotSize    int
}

func (img *CaptchaImage) calculateSizes(width, height, ncount int) {
	// Goal: fit all numbers into the image.
	// Convert everything to floats for calculations.
	w := float64(width)
	h := float64(height)
	// fontWidth includes 1-dot spacing between numbers
	fw := float64(fontWidth) + 1
	fh := float64(fontHeight)
	nc := float64(ncount)
	// Calculate width of a sigle number if we only take into
	// account the width
	nw := w / nc
	// Calculate the number height from this width
	nh := nw * fh / fw
	// Number height too large?
	if nh > h {
		// Fit numbers based on height
		nh = h
		nw = fw/fh * nh
	}
	// Calculate dot size
	img.dotSize = int(nh / fh)
	// Save everything, making actual width smaller by 1 dot,
	// to account for spacing between numbers
	img.numberWidth = int(nw)
	img.numberHeight = int(nh) - img.dotSize
}

func NewImage(numbers []byte, width, height int) *CaptchaImage {
	img := new(CaptchaImage)
	img.NRGBA = image.NewNRGBA(width, height)
	img.primaryColor = image.NRGBAColor{uint8(rand.Intn(50)), uint8(rand.Intn(50)), uint8(rand.Intn(128)), 0xFF}
	// We need some space, so calculate border
	var border int
	if width > height {
		border = height/5
	} else {
		border = width/5
	}
	bwidth := width-border*2
	bheight := height-border*2
	img.calculateSizes(bwidth, bheight, len(numbers))
	// Center numbers in image
	x := width/2 - (img.numberWidth*len(numbers))/2 - img.dotSize
	y := height/2 - img.numberHeight/2 + img.dotSize/2
	setRandomBrightness(&img.primaryColor, 180)
	for _, n := range numbers {
		//y = rand.Intn(dotSize * 4)
		img.drawNumber(font[n], x, y)
		x += img.numberWidth + img.dotSize
	}
	//img.strikeThrough(img.primaryColor)
	return img
}

func NewRandomImage(width, height int) (img *CaptchaImage, numbers []byte) {
	numbers = randomNumbers()
	img = NewImage(numbers, width, height)
	return
}

func (img *CaptchaImage) PNGEncode(w io.Writer) os.Error {
	return png.Encode(w, img)
}

func (img *CaptchaImage) drawHorizLine(color image.Color, fromX, toX, y int) {
	for x := fromX; x <= toX; x++ {
		img.Set(x, y, color)
	}
}

func (img *CaptchaImage) drawCircle(color image.Color, x, y, radius int) {
	f := 1 - radius
	dfx := 1
	dfy := -2 * radius
	xx := 0
	yy := radius

	img.Set(x, y+radius, color)
	img.Set(x, y-radius, color)
	img.drawHorizLine(color, x-radius, x+radius, y)

	for xx < yy {
		if f >= 0 {
			yy--
			dfy += 2
			f += dfy
		}
		xx++
		dfx += 2
		f += dfx
		img.drawHorizLine(color, x-xx, x+xx, y+yy)
		img.drawHorizLine(color, x-xx, x+xx, y-yy)
		img.drawHorizLine(color, x-yy, x+yy, y+xx)
		img.drawHorizLine(color, x-yy, x+yy, y-xx)
	}
}


func min3(x, y, z uint8) (o uint8) {
	o = x
	if y < o {
		o = y
	}
	if z < o {
		o = z
	}
	return
}

func max3(x, y, z uint8) (o uint8) {
	o = x
	if y > o {
		o = y
	}
	if z > o {
		o = z
	}
	return
}

func setRandomBrightness(c *image.NRGBAColor, max uint8) {
	minc := min3(c.R, c.G, c.B)
	maxc := max3(c.R, c.G, c.B)
	if maxc > max {
		return
	}
	n := rand.Intn(int(max-maxc)) - int(minc)
	c.R = uint8(int(c.R) + n)
	c.G = uint8(int(c.G) + n)
	c.B = uint8(int(c.B) + n)
}

func rnd(from, to int) int {
	return rand.Intn(to+1-from) + from
}

func (img *CaptchaImage) fillWithCircles(color image.NRGBAColor, n, maxradius int) {
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	for i := 0; i < n; i++ {
		setRandomBrightness(&color, 255)
		r := rnd(1, maxradius)
		img.drawCircle(color, rnd(r, maxx), rnd(r, maxy), r)
	}
}

// func (img *CaptchaImage) strikeThrough(color image.Color) {
// 	r := 0
// 	maxx := img.Bounds().Max.X
// 	maxy := img.Bounds().Max.Y
// 	y := rnd(maxy/3, maxy-maxy/3)
// 	for x := 0; x < maxx; x += r {
// 		r = rnd(1, dotSize/2-1)
// 		y += rnd(-2, 2)
// 		if y <= 0 || y >= maxy {
// 			y = rnd(maxy/3, maxy-maxy/3)
// 		}
// 		img.drawCircle(color, x, y, r)
// 	}
// }

func (img *CaptchaImage) drawNumber(number []byte, x, y int) {
	//skf := rand.Intn(maxSkew) - maxSkew/2
	//if skf < 0 {
	//	x -= skf * numberHeight
	//}
	//d := img.numberWidth / fontWidth // number height is ignored
	//println(img.numberWidth)
	minr := img.dotSize/2 // minumum radius
	maxr := img.dotSize // maximum radius
	// x += srad
	// y += srad
	for yy := 0; yy < fontHeight; yy++ {
		for xx := 0; xx < fontWidth; xx++ {
			if number[yy*fontWidth+xx] != 1 {
				continue
			}
			// introduce random variations
			or := rnd(minr, maxr)
			ox := x + (xx * maxr) //+ rnd(0, or/2)
			oy := y + (yy * maxr) //+ rnd(0, or/2)
			img.drawCircle(img.primaryColor, ox, oy, or)
		}
		//x += skf
	}
}
