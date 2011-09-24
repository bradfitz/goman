include $(GOROOT)/src/Make.inc

TARG=danga/gearman
GOFILES=\
	gearman.go\
	client.go\
	call.go\
	worker.go\
	util.go\

include $(GOROOT)/src/Make.pkg

w1.6: w1.go
w1: w1.6
