// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// generate is a tool to generate sounds.go from WAVE files.
//
// It creates (or rewrites) sounds.go in the parent directory.
package main

import (
	"fmt"
	"os"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
)

const headerLen = 44

var langs = []string{"en", "ru", "zh", "ja"}

func writeVar(w io.Writer, b []byte, prefix string) {
	i := 0
	for j, v := range b {
		fmt.Fprintf(w, "0x%02x,", v)
		i++
		if i == 11 {
			fmt.Fprintf(w, "\n")
			if j != len(b)-1 {
				fmt.Fprintf(w, prefix)
			}
			i = 0
		} else {
			if j != len(b)-1 {
				fmt.Fprintf(w, " ")
			}
		}
	}
	if i > 0 {
		fmt.Fprintf(w, "\n")
	}
}

func writeFileRep(pcm io.Writer, name, prefix string) {
	b, err := ioutil.ReadFile(name)
	if err != nil {
		log.Fatalf("%s", err)
	}
	writeVar(pcm, b[headerLen:], prefix)
}

func writeSingle(pcm io.Writer, name string) {
	fmt.Fprintf(pcm, "\nvar %sSound = []byte{\n\t", name)
	writeFileRep(pcm, name+".wav", "\t")
	fmt.Fprintf(pcm, "}\n")
}

func writeDigitSounds(pcm io.Writer, lang string) {
	fmt.Fprintf(pcm, "\t\"%s\": [][]byte{\n", lang)
	for i := 0; i <= 9; i++ {
		fmt.Fprintf(pcm, "\t\t{ // %d\n\t\t\t", i)
		writeFileRep(pcm, filepath.Join(lang, fmt.Sprintf("%d.wav", i)), "\t\t\t")
		fmt.Fprintf(pcm, "\t\t},\n")
	}
	fmt.Fprintf(pcm, "\t},\n")
}

func main() {
	pcm, err := os.Create(filepath.Join("..", "sounds.go"))
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer pcm.Close()
	fmt.Fprintf(pcm, `package captcha

// This file has been generated from .wav files using generate.go.

var waveHeader = []byte{
	0x52, 0x49, 0x46, 0x46, 0x00, 0x00, 0x00, 0x00, 0x57, 0x41, 0x56, 0x45,
	0x66, 0x6d, 0x74, 0x20, 0x10, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00,
	0x40, 0x1f, 0x00, 0x00, 0x40, 0x1f, 0x00, 0x00, 0x01, 0x00, 0x08, 0x00,
	0x64, 0x61, 0x74, 0x61,
}

// Byte slices contain raw 8 kHz unsigned 8-bit PCM data (without wav header).

`)
	fmt.Fprintf(pcm, "var digitSounds = map[string][][]byte{\n")
	for _, lang := range langs {
		writeDigitSounds(pcm, lang)
	}
	fmt.Fprintf(pcm, "}\n")
	writeSingle(pcm, "beep")
}
