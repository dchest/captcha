package main

import (
	"github.com/dchest/captcha"
	"os"
)

func main() {
	img, _ := captcha.NewRandomImage(300, 80)
	img.PNGEncode(os.Stdout)
}
