// Copyright 2011 Dmitry Chestnykh. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"net/http"
	"path"
	"strconv"
	"strings"
)

type captchaHandler struct {
	imgWidth  int
	imgHeight int
}

// Server returns a handler that serves HTTP requests with image or
// audio representations of captchas. Image dimensions are accepted as
// arguments. The server decides which captcha to serve based on the last URL
// path component: file name part must contain a captcha id, file extension —
// its format (PNG or WAV).
//
// For example, for file name "LBm5vMjHDtdUfaWYXiQX.png" it serves an image captcha
// with id "LBm5vMjHDtdUfaWYXiQX", and for "LBm5vMjHDtdUfaWYXiQX.wav" it serves the
// same captcha in audio format.
//
// To serve a captcha as a downloadable file, the URL must be constructed in
// such a way as if the file to serve is in the "download" subdirectory:
// "/download/LBm5vMjHDtdUfaWYXiQX.wav".
//
// To reload captcha (get a different solution for the same captcha id), append
// "?reload=x" to URL, where x may be anything (for example, current time or a
// random number to make browsers refetch an image instead of loading it from
// cache).
//
// By default, the Server serves audio in English language. To serve audio
// captcha in one of the other supported languages, append "lang" value, for
// example, "?lang=ru".
func Server(imgWidth, imgHeight int) http.Handler {
	return &captchaHandler{imgWidth, imgHeight}
}

func (h *captchaHandler) serve(w http.ResponseWriter, id, ext string, lang string, download bool) error {
	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	switch ext {
	case ".png":
		if !download {
			w.Header().Set("Content-Type", "image/png")
		}
		return WriteImage(w, id, h.imgWidth, h.imgHeight)
	case ".wav":
		//XXX(dchest) Workaround for Chrome: it wants content-length,
		//or else will start playing NOT from the beginning.
		//Filed issue: http://code.google.com/p/chromium/issues/detail?id=80565
		d := globalStore.Get(id, false)
		if d == nil {
			return ErrNotFound
		}
		a := NewAudio(d, lang)
		if !download {
			w.Header().Set("Content-Type", "audio/x-wav")
		}
		w.Header().Set("Content-Length", strconv.Itoa(a.EncodedLen()))
		_, err := a.WriteTo(w)
		return err
	}
	return ErrNotFound
}

func (h *captchaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir, file := path.Split(r.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || id == "" {
		http.NotFound(w, r)
		return
	}
	if r.FormValue("reload") != "" {
		Reload(id)
	}
	lang := strings.ToLower(r.FormValue("lang"))
	download := path.Base(dir) == "download"
	if h.serve(w, id, ext, lang, download) == ErrNotFound {
		http.NotFound(w, r)
	}
	// Ignore other errors.
}
