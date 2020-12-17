package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/dchest/captcha"
	"github.com/go-redis/redis"
)

var formTemplate = template.Must(template.New("example").Parse(formTemplateSrc))

func showFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	d := struct {
		CaptchaId string
	}{
		captcha.New(),
	}
	if err := formTemplate.Execute(w, &d); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func processFormHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if !captcha.VerifyString(r.FormValue("captchaId"), r.FormValue("captchaSolution")) {
		io.WriteString(w, "Wrong captcha solution! No robots allowed!\n")
	} else {
		io.WriteString(w, "Great job, human! You solved the captcha.\n")
	}
	io.WriteString(w, "<br><a href='/'>Try another one</a>")
}

func main() {
	// redis store
	s, err := captcha.NewRedisStore(&redis.Options{Addr: "localhost:6379", DB: 0}, captcha.Expiration, captcha.DefaultMaxRedisKeys, captcha.DefaultRedisPrefixKey)
	if err != nil {
		panic(err.Error())
	}
	captcha.SetCustomStore(s)

	// http
	http.HandleFunc("/", showFormHandler)
	http.HandleFunc("/process", processFormHandler)
	http.Handle("/captcha/", captcha.Server(captcha.StdWidth, captcha.StdHeight))
	fmt.Println("Server is at localhost:8666")
	if err := http.ListenAndServe("localhost:8666", nil); err != nil {
		log.Fatal(err)
	}
}

const formTemplateSrc = `<!doctype html>
<head><title>Captcha Example</title></head>
<body>
<script>
function setSrcQuery(e, q) {
	var src  = e.src;
	var p = src.indexOf('?');
	if (p >= 0) {
		src = src.substr(0, p);
	}
	e.src = src + "?" + q
}

function playAudio() {
	var le = document.getElementById("lang");
	var lang = le.options[le.selectedIndex].value;
	var e = document.getElementById('audio')
	setSrcQuery(e, "lang=" + lang)
	e.style.display = 'block';
	e.autoplay = 'true';
	return false;
}

function changeLang() {
	var e = document.getElementById('audio')
	if (e.style.display == 'block') {
		playAudio();
	}
}

function reload() {
	setSrcQuery(document.getElementById('image'), "reload=" + (new Date()).getTime());
	setSrcQuery(document.getElementById('audio'), (new Date()).getTime());
	return false;
}
</script>
<select id="lang" onchange="changeLang()">
	<option value="en">English</option>
	<option value="ja">Japanese</option>
	<option value="ru">Russian</option>
	<option value="zh">Chinese</option>
</select>
<form action="/process" method=post>
<p>Type the numbers you see in the picture below:</p>
<p><img id=image src="/captcha/{{.CaptchaId}}.png" alt="Captcha image"></p>
<a href="#" onclick="reload()">Reload</a> | <a href="#" onclick="playAudio()">Play Audio</a>
<audio id=audio controls style="display:none" src="/captcha/{{.CaptchaId}}.wav" preload=none>
  You browser doesn't support audio.
  <a href="/captcha/download/{{.CaptchaId}}.wav">Download file</a> to play it in the external player.
</audio>
<input type=hidden name=captchaId value="{{.CaptchaId}}"><br>
<input name=captchaSolution>
<input type=submit value=Submit>
</form>
`
