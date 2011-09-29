O=6
GC=$(O)g
LD=$(O)l
.SUFFIXES : .go .$(O)

.go.6:
	$(GC) $<

.6:
	$(LD) -o $@ $<

include $(GOROOT)/src/Make.inc

TARG=danga/gearman
GOFILES=\
	gearman.go\
	client.go\
	call.go\
	worker.go\
	util.go\

include $(GOROOT)/src/Make.pkg
