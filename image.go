package captcha

import (
	"image"
	"image/png"
	"io"
	"os"
	"rand"
)

const (
	// Standard width and height for captcha image
	StdWidth  = 300
	StdHeight = 80

	maxSkew = 2
)

type CaptchaImage struct {
	*image.NRGBA
	primaryColor image.NRGBAColor
	numberWidth  int
	numberHeight int
	dotSize      int
}

// NewImage returns a new captcha image of the given width and height with the
// given slice of numbers, where each number must be in range 0-9.
func NewImage(numbers []byte, width, height int) *CaptchaImage {
	img := new(CaptchaImage)
	img.NRGBA = image.NewNRGBA(width, height)
	img.primaryColor = image.NRGBAColor{uint8(rand.Intn(50)), uint8(rand.Intn(50)), uint8(rand.Intn(128)), 0xFF}
	// We need some space, so calculate border
	var border int
	if width > height {
		border = height / 5
	} else {
		border = width / 5
	}
	bwidth := width - border*2
	bheight := height - border*2
	img.calculateSizes(bwidth, bheight, len(numbers))
	// Background 
	img.fillWithCircles(10, img.dotSize)
	maxx := width - (img.numberWidth+img.dotSize)*len(numbers) - img.dotSize
	maxy := height - img.numberHeight - img.dotSize*2
	x := rnd(img.dotSize*2, maxx)
	y := rnd(img.dotSize*2, maxy)
	setRandomBrightness(&img.primaryColor, 180)
	for _, n := range numbers {
		img.drawNumber(font[n], x, y)
		x += img.numberWidth + img.dotSize
	}
	img.strikeThrough()
	return img
}

// NewRandomImage generates random numbers and returns a new captcha image of
// the given width and height with those numbers printed on it, and the numbers
// themselves.
func NewRandomImage(width, height int) (img *CaptchaImage, numbers []byte) {
	numbers = randomNumbers()
	img = NewImage(numbers, width, height)
	return
}

// PNGEncode writes captcha image in PNG format into the given writer.
func (img *CaptchaImage) PNGEncode(w io.Writer) os.Error {
	return png.Encode(w, img)
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
		nw = fw / fh * nh
	}
	// Calculate dot size
	img.dotSize = int(nh / fh)
	// Save everything, making actual width smaller by 1 dot,
	// to account for spacing between numbers
	img.numberWidth = int(nw)
	img.numberHeight = int(nh) - img.dotSize
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

func (img *CaptchaImage) fillWithCircles(n, maxradius int) {
	color := img.primaryColor
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	for i := 0; i < n; i++ {
		setRandomBrightness(&color, 255)
		r := rnd(1, maxradius)
		img.drawCircle(color, rnd(r, maxx-r), rnd(r, maxy-r), r)
	}
}

func (img *CaptchaImage) strikeThrough() {
	r := 0
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	y := rnd(maxy/3, maxy-maxy/3)
	for x := 0; x < maxx; x += r {
		r = rnd(1, img.dotSize/2-1)
		y += rnd(-img.dotSize/2, img.dotSize/2)
		if y <= 0 || y >= maxy {
			y = rnd(maxy/3, maxy-maxy/3)
		}
		img.drawCircle(img.primaryColor, x, y, r)
	}
}

func (img *CaptchaImage) drawNumber(number []byte, x, y int) {
	skf := rand.Float64() * float64(rnd(-maxSkew, maxSkew))
	xs := float64(x)
	minr := img.dotSize / 2               // minumum radius
	maxr := img.dotSize/2 + img.dotSize/4 // maximum radius
	y += rnd(-minr, minr)
	for yy := 0; yy < fontHeight; yy++ {
		for xx := 0; xx < fontWidth; xx++ {
			if number[yy*fontWidth+xx] != 1 {
				continue
			}
			// introduce random variations
			or := rnd(minr, maxr)
			ox := x + (xx * img.dotSize) + rnd(0, or/2)
			oy := y + (yy * img.dotSize) + rnd(0, or/2)
			img.drawCircle(img.primaryColor, ox, oy, or)
		}
		xs += skf
		x = int(xs)
	}
}
