include $(GOROOT)/src/Make.inc

TARG=github.com/dchest/captcha
GOFILES=\
	captcha.go\
	draw.go\
	font.go

include $(GOROOT)/src/Make.pkg

