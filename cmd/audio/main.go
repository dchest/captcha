package main

import (
	"github.com/dchest/captcha"
	"os"
)

func main() {
	c, _ := captcha.NewRandomAudio(captcha.StdLength)
	c.WriteTo(os.Stdout)
}
