package captcha

import (
	"image"
	"image/png"
	"io"
	"math"
	"os"
	"rand"
	"time"
)

const (
	// Standard width and height of a captcha image.
	StdWidth  = 240
	StdHeight = 80
	// Maximum absolute skew factor of a single digit.
	maxSkew = 0.7
	// Number of background circles.
	circleCount = 20
)

type Image struct {
	*image.Paletted
	numWidth  int
	numHeight int
	dotSize   int
}

func init() {
	rand.Seed(time.Seconds())
}

func randomPalette() image.PalettedColorModel {
	p := make([]image.Color, circleCount+1)
	// Transparent color.
	// TODO(dchest). Currently it's white, not transparent, because PNG
	// encoder doesn't support paletted images with alpha channel.
	// Submitted CL: http://codereview.appspot.com/4432078 Change alpha to
	// 0x00 once it's accepted.
	p[0] = image.RGBAColor{0xFF, 0xFF, 0xFF, 0xFF}
	// Primary color.
	prim := image.RGBAColor{
		uint8(rand.Intn(129)),
		uint8(rand.Intn(129)),
		uint8(rand.Intn(129)),
		0xFF,
	}
	p[1] = prim
	// Circle colors.
	for i := 2; i <= circleCount; i++ {
		p[i] = randomBrightness(prim, 255)
	}
	return p
}

// NewImage returns a new captcha image of the given width and height with the
// given digits, where each digit must be in range 0-9.
func NewImage(digits []byte, width, height int) *Image {
	img := new(Image)
	img.Paletted = image.NewPaletted(width, height, randomPalette())
	img.calculateSizes(width, height, len(digits))
	// Randomly position captcha inside the image.
	maxx := width - (img.numWidth+img.dotSize)*len(digits) - img.dotSize
	maxy := height - img.numHeight - img.dotSize*2
	var border int
	if width > height {
		border = height / 5
	} else {
		border = width / 5
	}
	x := rnd(border, maxx-border)
	y := rnd(border, maxy-border)
	// Draw digits.
	for _, n := range digits {
		img.drawDigit(font[n], x, y)
		x += img.numWidth + img.dotSize
	}
	// Draw strike-through line.
	img.strikeThrough()
	// Apply wave distortion.
	img.distort(rndf(5, 10), rndf(100, 200))
	// Fill image with random circles.
	img.fillWithCircles(circleCount, img.dotSize)
	return img
}

// BUG(dchest): While Image conforms to io.WriterTo interface, its WriteTo
// method returns 0 instead of the actual bytes written because png.Encode
// doesn't report this.

// WriteTo writes captcha image in PNG format into the given writer.
func (img *Image) WriteTo(w io.Writer) (int64, os.Error) {
	return 0, png.Encode(w, img.Paletted)
}

func (img *Image) calculateSizes(width, height, ncount int) {
	// Goal: fit all digits inside the image.
	var border int
	if width > height {
		border = height / 4
	} else {
		border = width / 4
	}
	// Convert everything to floats for calculations.
	w := float64(width - border*2)
	h := float64(height - border*2)
	// fw takes into account 1-dot spacing between digits.
	fw := float64(fontWidth + 1)
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
	img.numWidth = int(nw) - img.dotSize
	img.numHeight = int(nh)
}

func (img *Image) drawHorizLine(fromX, toX, y int, colorIndex uint8) {
	for x := fromX; x <= toX; x++ {
		img.SetColorIndex(x, y, colorIndex)
	}
}

func (img *Image) drawCircle(x, y, radius int, colorIndex uint8) {
	f := 1 - radius
	dfx := 1
	dfy := -2 * radius
	xx := 0
	yy := radius

	img.SetColorIndex(x, y+radius, colorIndex)
	img.SetColorIndex(x, y-radius, colorIndex)
	img.drawHorizLine(x-radius, x+radius, y, colorIndex)

	for xx < yy {
		if f >= 0 {
			yy--
			dfy += 2
			f += dfy
		}
		xx++
		dfx += 2
		f += dfx
		img.drawHorizLine(x-xx, x+xx, y+yy, colorIndex)
		img.drawHorizLine(x-xx, x+xx, y-yy, colorIndex)
		img.drawHorizLine(x-yy, x+yy, y+xx, colorIndex)
		img.drawHorizLine(x-yy, x+yy, y-xx, colorIndex)
	}
}

func (img *Image) fillWithCircles(n, maxradius int) {
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	for i := 0; i < n; i++ {
		colorIndex := uint8(rnd(1, circleCount-1))
		r := rnd(1, maxradius)
		img.drawCircle(rnd(r, maxx-r), rnd(r, maxy-r), r, colorIndex)
	}
}

func (img *Image) strikeThrough() {
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	y := rnd(maxy/3, maxy-maxy/3)
	amplitude := rndf(5, 20)
	period := rndf(80, 180)
	dx := 2.0 * math.Pi / period
	for x := 0; x < maxx; x++ {
		xo := amplitude * math.Cos(float64(y)*dx)
		yo := amplitude * math.Sin(float64(x)*dx)
		for yn := 0; yn < img.dotSize; yn++ {
			r := rnd(0, img.dotSize)
			img.drawCircle(x+int(xo), y+int(yo)+(yn*img.dotSize), r/2, 1)
		}
	}
}

func (img *Image) drawDigit(digit []byte, x, y int) {
	skf := rndf(-maxSkew, maxSkew)
	xs := float64(x)
	r := img.dotSize / 2
	y += rnd(-r, r)
	for yy := 0; yy < fontHeight; yy++ {
		for xx := 0; xx < fontWidth; xx++ {
			if digit[yy*fontWidth+xx] != blackChar {
				continue
			}
			ox := x + xx*img.dotSize
			oy := y + yy*img.dotSize
			img.drawCircle(ox, oy, r, 1)
		}
		xs += skf
		x = int(xs)
	}
}

func (img *Image) distort(amplude float64, period float64) {
	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	oldImg := img.Paletted
	newImg := image.NewPaletted(w, h, oldImg.Palette)

	dx := 2.0 * math.Pi / period
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			ox := amplude * math.Sin(float64(y)*dx)
			oy := amplude * math.Cos(float64(x)*dx)
			newImg.SetColorIndex(x, y, oldImg.ColorIndexAt(x+int(ox), y+int(oy)))
		}
	}
	img.Paletted = newImg
}

func randomBrightness(c image.RGBAColor, max uint8) image.RGBAColor {
	minc := min3(c.R, c.G, c.B)
	maxc := max3(c.R, c.G, c.B)
	if maxc > max {
		return c
	}
	n := rand.Intn(int(max-maxc)) - int(minc)
	return image.RGBAColor{
		uint8(int(c.R) + n),
		uint8(int(c.G) + n),
		uint8(int(c.B) + n),
		uint8(c.A),
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

// rnd returns a random number in range [from, to].
func rnd(from, to int) int {
	return rand.Intn(to+1-from) + from
}
