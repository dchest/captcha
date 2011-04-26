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
	fmt.Fprintf(w, formJs)
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

const formJs = `
<script>
function playAudio() {
	var e = document.getElementById('audio')
	e.style.display = 'block';
	e.play();
	return false;
}

function reload() {
	function setSrcQuery(e, q) {
		var src  = e.src;
		var p = src.indexOf('?');
		if (p >= 0) {
			src = src.substr(0, p);
		}
		e.src = src + "?" + q
	}
	setSrcQuery(document.getElementById('image'), "reload=" + (new Date()).getTime());
	setSrcQuery(document.getElementById('audio'), (new Date()).getTime());
	return false;
}
</script>
`

const formTemplateSrc = `
<form action="/process" method=post>
<p>Type the numbers you see in the picture below:</p>
<p><img id=image src="/captcha/{CaptchaId}.png" alt="Captcha image"></p>
<a href="#" onclick="return reload()">Reload</a> | <a href="#" onclick="return playAudio()">Play Audio</a>
<audio id=audio controls style="display:none" src="/captcha/{CaptchaId}.wav" preload=none type="audio/wav">
  You browser doesn't support audio.
  <a href="/captcha/{CaptchaId}.wav?get">Download file</a> to play it in the external player.
</audio>
<input type=hidden name=captchaId value={CaptchaId}><br>
<input name=captchaSolution>
<input type=submit value=Submit>
`
