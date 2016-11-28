# (c) 2015 Wojciech A. Koszek <wojciech@koszek.com>

SRCS:=$(wildcard *.go)
PROGS:=$(SRCS:.go=.prog)

objs: $(PROGS)

%.prog: %.go
	go build -i -o $@ $<

clean:
	rm -rf *.prog
