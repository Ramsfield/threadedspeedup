GXX=gcc
CFLAGS= -std=c17 -Wall --pedantic
CLIBS= -pthread

.PHONY: all clean test
all: rtest slowrtest

rtest: randtest.c
	${GXX} randtest.c -o rtest ${CFLAGS} ${CLIBS}

slowrtest: randtest.c
	${GXX} randtest.c -o slowrtest -DSLOWDOWN ${CFLAGS} ${CLIBS}

clean:
	rm -rf slowrtest rtest
