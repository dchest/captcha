package captcha

import (
	"http"
	"os"
	"path"
	"strconv"
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
// For example, for file name "B9QTvDV1RXbVJ3Ac.png" it serves an image captcha
// with id "B9QTvDV1RXbVJ3Ac", and for "B9QTvDV1RXbVJ3Ac.wav" it serves the
// same captcha in audio format.
//
// To serve an audio captcha as downloadable file, append "?get" to URL.
//
// To reload captcha (get a different solution for the same captcha id), append
// "?reload=x" to URL, where x may be anything (for example, current time or a
// random number to make browsers refetch an image instead of loading it from
// cache).
func Server(w, h int) http.Handler { return &captchaHandler{w, h} }

func (h *captchaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, file := path.Split(r.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || id == "" {
		http.NotFound(w, r)
		return
	}
	if r.FormValue("reload") != "" {
		Reload(id)
	}
	var err os.Error
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		err = WriteImage(w, id, h.imgWidth, h.imgHeight)
	case ".wav":
		if r.URL.RawQuery == "get" {
			w.Header().Set("Content-Type", "application/octet-stream")
		} else {
			w.Header().Set("Content-Type", "audio/x-wav")
		}
		//err = WriteAudio(w, id)
		//XXX(dchest) Workaround for Chrome: it wants content-length,
		//or else will start playing NOT from the beginning.
		//Filed issue: http://code.google.com/p/chromium/issues/detail?id=80565
		d := globalStore.Get(id, false)
		if d == nil {
			err = ErrNotFound
		} else {
			a := NewAudio(d)
			w.Header().Set("Content-Length", strconv.Itoa(a.EncodedLen()))
			_, err = a.WriteTo(w)
		}
	default:
		err = ErrNotFound
	}
	if err != nil {
		if err == ErrNotFound {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "error serving captcha", http.StatusInternalServerError)
	}
}
