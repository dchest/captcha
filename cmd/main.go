package main

import (
	"github.com/dchest/captcha"
	"os"
)

func main() {
	img, _ := captcha.NewRandomImage(captcha.StdWidth, captcha.StdHeight)
	img.PNGEncode(os.Stdout)
}
