// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package captcha

import (
	"bytes"
	"github.com/valyala/fasthttp"
	"io"
	"net/http"
	"path"
	"strings"
)

type captchaFastHTTPHandler struct {
	imgWidth  int
	imgHeight int
}

func ServerFastHTTP(imgWidth, imgHeight int) *captchaFastHTTPHandler {
	return &captchaFastHTTPHandler{imgWidth, imgHeight}
}

func (h *captchaFastHTTPHandler) serveFastHTTP(ctx *fasthttp.RequestCtx, id, ext, lang string, download bool) error {
	ctx.Response.Header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Response.Header.Set("Pragma", "no-cache")
	ctx.Response.Header.Set("Expires", "0")

	var content bytes.Buffer
	switch ext {
	case ".png":
		ctx.Response.Header.Set("Content-Type", "image/png")
		WriteImage(&content, id, h.imgWidth, h.imgHeight)
	case ".wav":
		ctx.Response.Header.Set("Content-Type", "audio/x-wav")
		WriteAudio(&content, id, lang)
	default:
		return ErrNotFound
	}

	if download {
		ctx.Response.Header.Set("Content-Type", "application/octet-stream")
	}

	ctx.SetStatusCode(http.StatusOK)
	var out io.Writer
	out = ctx.Response.BodyWriter()
	io.Copy(out, &content)

	return nil
}

func (h *captchaFastHTTPHandler) ServeFastHTTP(ctx *fasthttp.RequestCtx) {
	dir, file := path.Split(string(ctx.URI().Path()))
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || id == "" {
		ctx.NotFound()
		return
	}
	if len(ctx.FormValue("reload")) > 0 {
		Reload(id)
	}
	lang := strings.ToLower(string(ctx.FormValue("lang")))
	download := path.Base(dir) == "download"
	if h.serveFastHTTP(ctx, id, ext, lang, download) == ErrNotFound {
		ctx.NotFound()
	}
}

func VerifyBytes(id string, digits []byte) bool {
	if len(digits) == 0 {
		return false
	}
	ns := make([]byte, len(digits))
	for i := range ns {
		d := digits[i]
		switch {
		case '0' <= d && d <= '9':
			ns[i] = d - '0'
		case d == ' ' || d == ',':
		// ignore
		default:
			return false
		}
	}
	return Verify(id, ns)
}
