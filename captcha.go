package main

import (
	"image"
	"image/png"
	"os"
	"rand"
	"time"
	crand "crypto/rand"
	"io"
)

var numbers = [][]byte{
	{
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 0, 0, 0, 1,
		0, 1, 1, 1, 0,
	},
	{
		0, 0, 1, 0, 0,
		0, 1, 1, 0, 0,
		1, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		0, 0, 1, 0, 0,
		1, 1, 1, 1, 1,
	},
	{
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 1, 1,
		0, 1, 1, 0, 0,
		1, 0, 0, 0, 0,
		1, 0, 0, 0, 0,
		1, 1, 1, 1, 1,
	},
	{
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 1, 1,
		0, 1, 1, 0, 0,
		0, 0, 0, 1, 0,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		1, 1, 1, 1, 0,
	},
	{
		1, 0, 0, 1, 0,
		1, 0, 0, 1, 0,
		1, 0, 0, 1, 0,
		1, 0, 0, 1, 0,
		1, 1, 1, 1, 1,
		0, 0, 0, 1, 0,
		0, 0, 0, 1, 0,
		0, 0, 0, 1, 0,
	},
	{
		1, 1, 1, 1, 1,
		1, 0, 0, 0, 0,
		1, 0, 0, 0, 0,
		1, 1, 1, 1, 0,
		0, 0, 0, 1, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 1, 1,
		1, 1, 1, 1, 0,
	},
	{
		0, 0, 1, 1, 1,
		0, 1, 0, 0, 0,
		1, 0, 0, 0, 0,
		1, 1, 1, 1, 0,
		1, 1, 0, 0, 1,
		1, 0, 0, 0, 1,
		1, 1, 0, 0, 1,
		0, 1, 1, 1, 0,
	},
	{
		1, 1, 1, 1, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 1, 0,
		0, 0, 1, 0, 0,
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
		0, 1, 0, 0, 0,
	},
	{
		0, 1, 1, 1, 0,
		1, 0, 0, 0, 1,
		1, 1, 0, 1, 1,
		0, 1, 1, 1, 0,
		1, 1, 0, 1, 1,
		1, 0, 0, 0, 1,
		1, 1, 0, 1, 1,
		0, 1, 1, 1, 0,
	},
	{
		0, 1, 1, 1, 0,
		1, 0, 0, 1, 1,
		1, 0, 0, 0, 1,
		1, 1, 0, 0, 1,
		0, 1, 1, 1, 1,
		0, 0, 0, 0, 1,
		0, 0, 0, 0, 1,
		1, 1, 1, 1, 0,
	},
}

const (
	NumberWidth  = 5
	NumberHeight = 8
	DotSize = 6
	SkewFactor = 3
)

func drawHorizLine(img *image.NRGBA, color image.Color, fromX, toX, y int) {
	for x := fromX; x <= toX; x++ {
		img.Set(x, y, color)
	}
}

func drawCircle(img *image.NRGBA, color image.Color, x0, y0, radius int) {
	f := 1 - radius
	ddF_x := 1
	ddF_y := -2 * radius
	x := 0
	y := radius

	img.Set(x0, y0+radius, color)
	img.Set(x0, y0-radius, color)
	//img.Set(x0+radius, y0, color)
	//img.Set(x0-radius, y0, color)
	drawHorizLine(img, color, x0-radius, x0+radius, y0)

	for x < y {
		// ddF_x == 2 * x + 1;
		// ddF_y == -2 * y;
		// f == x*x + y*y - radius*radius + 2*x - y + 1;
		if f >= 0 {
			y--
			ddF_y += 2
			f += ddF_y
		}
		x++
		ddF_x += 2
		f += ddF_x
		//img.Set(x0+x, y0+y, color)
		//img.Set(x0-x, y0+y, color)
		drawHorizLine(img, color, x0-x, x0+x, y0+y)
		//img.Set(x0+x, y0-y, color)
		//img.Set(x0-x, y0-y, color)
		drawHorizLine(img, color, x0-x, x0+x, y0-y)
		//img.Set(x0+y, y0+x, color)
		//img.Set(x0-y, y0+x, color)
		drawHorizLine(img, color, x0-y, x0+y, y0+x)
		//img.Set(x0+y, y0-x, color)
		//img.Set(x0-y, y0-x, color)
		drawHorizLine(img, color, x0-y, x0+y, y0-x)
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

func fillWithCircles(img *image.NRGBA, n, maxradius int) {
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	color := image.NRGBAColor{0, 0, 0x80, 0xFF}
	for i := 0; i < n; i++ {
		setRandomBrightness(&color, 255)
		r := rand.Intn(maxradius-1)+1 
		drawCircle(img, color, rand.Intn(maxx-r*2)+r, rand.Intn(maxy-r*2)+r, r)
	}
}

func drawNumber(img *image.NRGBA, number []byte, x, y int, color image.NRGBAColor) {
	skf := rand.Intn(SkewFactor)-SkewFactor/2
	if skf < 0 {
		x -= skf * NumberHeight
	}
	for y0 := 0; y0 < NumberHeight; y0++ {
		for x0 := 0; x0 < NumberWidth; x0++ {
			radius := rand.Intn(DotSize/2)+DotSize/2
			addx := rand.Intn(radius/2)
			addy := rand.Intn(radius/2)
			if number[y0*NumberWidth+x0] == 1 {
				drawCircle(img, color, x+x0*DotSize+DotSize+addx, y+y0*DotSize+DotSize+addy, radius)
			}
		}
		x += skf
	}
}

func drawNumbersToImage(ns []byte) {
	img := image.NewNRGBA(NumberWidth*(DotSize+2)*len(ns)+DotSize, NumberHeight*DotSize+(DotSize*6))
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, image.NRGBAColor{0xFF, 0xFF, 0xFF, 0xFF})
		}
	}
	fillWithCircles(img, 60, 3)
	x := rand.Intn(DotSize)
	y := 0
	color := image.NRGBAColor{0, 0, 0x80, 0xFF}
	setRandomBrightness(&color, 180)
	for _, n := range ns {
		y = rand.Intn(DotSize*4)
		drawNumber(img, numbers[n], x, y, color)
		x += DotSize * NumberWidth + rand.Intn(SkewFactor)+3 
	}
	f, err := os.Create("captcha.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	png.Encode(f, img)
}

func main() {
	rand.Seed(time.Seconds())
	n := make([]byte, 6)
	if _, err := io.ReadFull(crand.Reader, n); err != nil {
		panic(err)
	}
	for i := range n {
		n[i] %= 10
	}
	drawNumbersToImage(n)
}

