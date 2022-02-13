GXX=gcc
CFLAGS= -std=c17 -Wall --pedantic
CLIBS= -pthread

rtest: randtest.c
	${GXX} randtest.c -o rtest ${CFLAGS} ${CLIBS}

.PHONY: clean
clean:
	rm -rf rtest
