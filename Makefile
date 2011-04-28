include $(GOROOT)/src/Make.inc

TARG=github.com/dchest/captcha
GOFILES=\
	captcha.go\
	random.go\
	store.go\
	font.go\
	image.go\
	sounds.go\
	audio.go\
	server.go

include $(GOROOT)/src/Make.pkg

