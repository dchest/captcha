// example of HTTP server that uses the captcha package.
package main

import (
	"fmt"
	"github.com/dchest/captcha"
	"http"
	"template"
)

var formTemplate = template.MustParse(formTemplateSrc, nil)

func showFormHandler(w http.ResponseWriter, r *http.Request) {
	d := struct{ CaptchaId string }{captcha.New(captcha.StdLength)}
	if err := formTemplate.Execute(w, &d); err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
	}
}

func processFormHandler(w http.ResponseWriter, r *http.Request) {
	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		fmt.Fprintf(w, "Wrong captcha solution! No robots allowed!")
		return
	}
	fmt.Fprintf(w, "Great job, human! You solved the captcha.")
}

func main() {
	http.HandleFunc("/", showFormHandler)
	http.HandleFunc("/process", processFormHandler)
	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	http.ListenAndServe(":8080", nil)
}

const formTemplateSrc = `
<form action="/process" method=post>
<p>Type the numbers you see in the picture below:</p>
<p><img src="/captcha/{CaptchaId}.png" alt="Captcha image"></p>
<a href="#" onclick="var e=getElementById('audio'); e.display=true; e.play(); return false">Play Audio</a>
<audio id=audio controls style="display:none">
  <source src="/captcha/{CaptchaId}.wav" type="audio/x-wav">
  You browser doesn't support audio.
  <a href="/captcha/{CaptchaId}.wav?get">Download file</a> to play it in the external player.
</audio>
<input type=hidden name=captchaId value={CaptchaId}><br>
<input name=captchaSolution>
<input type=submit>
`
