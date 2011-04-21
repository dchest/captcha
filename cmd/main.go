package main

import (
	"github.com/dchest/captcha"
	"os"
)

func main() {
	img, _ := captcha.NewRandomImage(230, 60)
	img.PNGEncode(os.Stdout)
}
