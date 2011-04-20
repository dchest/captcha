package captcha

import (
	"image"
	"rand"
)

func drawHorizLine(img *image.NRGBA, color image.Color, fromX, toX, y int) {
	for x := fromX; x <= toX; x++ {
		img.Set(x, y, color)
	}
}

func drawCircle(img *image.NRGBA, color image.Color, x0, y0, radius int) {
	f := 1 - radius
	ddFx := 1
	ddFy := -2 * radius
	x := 0
	y := radius

	img.Set(x0, y0+radius, color)
	img.Set(x0, y0-radius, color)
	drawHorizLine(img, color, x0-radius, x0+radius, y0)

	for x < y {
		if f >= 0 {
			y--
			ddFy += 2
			f += ddFy
		}
		x++
		ddFx += 2
		f += ddFx
		drawHorizLine(img, color, x0-x, x0+x, y0+y)
		drawHorizLine(img, color, x0-x, x0+x, y0-y)
		drawHorizLine(img, color, x0-y, x0+y, y0+x)
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

func fillWithCircles(img *image.NRGBA, color image.NRGBAColor, n, maxradius int) {
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	for i := 0; i < n; i++ {
		setRandomBrightness(&color, 255)
		r := rand.Intn(maxradius-1) + 1
		drawCircle(img, color, rand.Intn(maxx-r*2)+r, rand.Intn(maxy-r*2)+r, r)
	}
}

func drawCirclesLine(img *image.NRGBA, color image.Color) {
	r := 0
	maxx := img.Bounds().Max.X
	maxy := img.Bounds().Max.Y
	y := rand.Intn(maxy/2+maxy/3) - rand.Intn(maxy/3)
	for x := 0; x < maxx; x += r  {
		r = rand.Intn(dotSize/2)
		y += rand.Intn(3) - 1
		if y <= 0 || y >= maxy {
			y = rand.Intn(maxy/2) + rand.Intn(maxy/2)
		}
		drawCircle(img, color, x, y, r)
	}
}

func drawNumber(img *image.NRGBA, number []byte, x, y int, color image.NRGBAColor) {
	skf := rand.Intn(maxSkew) - maxSkew/2
	if skf < 0 {
		x -= skf * numberHeight
	}
	for y0 := 0; y0 < numberHeight; y0++ {
		for x0 := 0; x0 < numberWidth; x0++ {
			radius := rand.Intn(dotSize/2) + dotSize/2
			addx := rand.Intn(radius / 4)
			addy := rand.Intn(radius / 4)
			if number[y0*numberWidth+x0] == 1 {
				drawCircle(img, color, x+x0*dotSize+dotSize+addx,
					y+y0*dotSize+dotSize+addy, radius)
			}
		}
		x += skf
	}
}
