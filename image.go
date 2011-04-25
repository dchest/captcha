package captcha

import (
	"image"
	"image/png"
	"io"
	"os"
	"rand"
	"time"
)

const (
	// Standard width and height of a captcha image.
	StdWidth  = 300
	StdHeight = 80

	maxSkew = 2
)

type Image struct {
	*image.NRGBA
	primaryColor image.NRGBAColor
	numWidth     int
	numHeight    int
	dotSize      int
}

func init() {
	rand.Seed(time.Seconds())
}

// NewImage returns a new captcha image of the given width and height with the
// given digits, where each digit must be in range 0-9.
func NewImage(digits []byte, width, height int) *Image {
	img := new(Image)
	img.NRGBA = image.NewNRGBA(width, height)
	img.primaryColor = image.NRGBAColor{
		uint8(rand.Intn(129)),
		uint8(rand.Intn(129)),
		uint8(rand.Intn(129)),
		0xFF,
	}
	img.calculateSizes(width, height, len(digits))
	// Draw background (10 random circles of random brightness).
	img.fillWithCircles(10, img.dotSize)
	// Randomly position captcha inside the image.
	maxx := width - (img.numWidth+img.dotSize)*len(digits) - img.dotSize
	maxy := height - img.numHeight - img.dotSize*2
	x := rnd(img.dotSize*2, maxx)
	y := rnd(img.dotSize*2, maxy)
	// Draw digits.
	for _, n := range digits {
		img.drawDigit(font[n], x, y)
		x += img.numWidth + img.dotSize
	}
	// Draw strike-through line.
	img.strikeThrough()
	return img
}

// BUG(dchest): While Image conforms to io.WriterTo interface, its WriteTo
// method returns 0 instead of the actual bytes written because png.Encode
// doesn't report this.

// WriteTo writes captcha image in PNG format into the given writer.
func (img *Image) WriteTo(w io.Writer) (int64, os.Error) {
	return 0, png.Encode(w, img)
}

func (img *Image) calculateSizes(width, height, ncount int) {
	// Goal: fit all digits inside the image.
	var border int
	if width > height {
		border = height / 5
	} else {
		border = width / 5
	}
	// Convert everything to floats for calculations.
	w := float64(width - border*2)
	h := float64(height - border*2)
	// fw takes into account 1-dot spacing between digits.
	fw := float64(fontWidth) + 1
	fh := float64(fontHeight)
	nc := float64(ncount)
	// Calculate the width of a single digit taking into account only the
	// width of the image.
	nw := w / nc
	// Calculate the height of a digit from this width.
	nh := nw * fh / fw
	// Digit too high?
	if nh > h {
		// Fit digits based on height.
		nh = h
		nw = fw / fh * nh
	}
	// Calculate dot size.
	img.dotSize = int(nh / fh)
	// Save everything, making the actual width smaller by 1 dot to account
	// for spacing between digits.
	img.numWidth = int(nw)
	img.numHeight = int(nh) - img.dotSize
}

func (img *Image) drawHorizLine(color image.Color, fromX, toX, y int) {
	for x := fromX; x <= toX; x++ {
		img.Set(x, y, color)
	}
}

func (img *Image) drawCircle(color image.Color, x, y, radius int) {
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

func (img *Image) fillWithCircles(n, maxradius int) {
	color := img.primaryColor
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	for i := 0; i < n; i++ {
		setRandomBrightness(&color, 255)
		r := rnd(1, maxradius)
		img.drawCircle(color, rnd(r, maxx-r), rnd(r, maxy-r), r)
	}
}

func (img *Image) strikeThrough() {
	r := 0
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	y := rnd(maxy/3, maxy-maxy/3)
	for x := 0; x < maxx; x += r {
		r = rnd(1, img.dotSize/3)
		y += rnd(-img.dotSize/2, img.dotSize/2)
		if y <= 0 || y >= maxy {
			y = rnd(maxy/3, maxy-maxy/3)
		}
		img.drawCircle(img.primaryColor, x, y, r)
	}
}

func (img *Image) drawDigit(digit []byte, x, y int) {
	skf := rand.Float64() * float64(rnd(-maxSkew, maxSkew))
	xs := float64(x)
	minr := img.dotSize / 2               // minumum radius
	maxr := img.dotSize/2 + img.dotSize/4 // maximum radius
	y += rnd(-minr, minr)
	for yy := 0; yy < fontHeight; yy++ {
		for xx := 0; xx < fontWidth; xx++ {
			if digit[yy*fontWidth+xx] != blackChar {
				continue
			}
			// Introduce random variations.
			or := rnd(minr, maxr)
			ox := x + (xx * img.dotSize) + rnd(0, or/2)
			oy := y + (yy * img.dotSize) + rnd(0, or/2)

			img.drawCircle(img.primaryColor, ox, oy, or)
		}
		xs += skf
		x = int(xs)
	}
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

// rnd returns a random number in range [from, to].
func rnd(from, to int) int {
	return rand.Intn(to+1-from) + from
}
