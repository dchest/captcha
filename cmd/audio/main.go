package main

import (
	"github.com/dchest/captcha"
	"log"
	"os"
)

func main() {
	f, err := os.Create("mixed.wav")
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer f.Close()

	c := captcha.NewAudio([]byte{1, 2, 3, 4, 5, 6})
	n, err := c.WriteTo(f)
	if err != nil {
		log.Fatalf("%s", err)
	}
	println("written", n)
}
