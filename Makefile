include $(GOROOT)/src/Make.inc

TARG=github.com/dchest/captcha
GOFILES=\
	captcha.go\
	store.go\
	font.go\
	image.go

include $(GOROOT)/src/Make.pkg

