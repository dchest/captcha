include $(GOROOT)/src/Make.inc

TARG=github.com/dchest/captcha
GOFILES=\
	captcha.go\
	store.go\
	font.go\
	image.go\
	sounds.go\
	audio.go

include $(GOROOT)/src/Make.pkg

