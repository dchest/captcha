package main

import (
	"github.com/dchest/captcha"
	"os"
)

func main() {
	captcha.Encode(os.Stdout)
}
