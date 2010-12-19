include $(GOROOT)/src/Make.inc

TARG=danga/gearman
GOFILES=\
	gearman.go\
	client.go\
	call.go\
	worker.go\

include $(GOROOT)/src/Make.pkg
