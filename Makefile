include $(GOROOT)/src/Make.inc

TARG=gocss
GOFILES=src/main/gocss.go
O_FILES=lexer.$O parser.$O

all: $(O_FILES)
install: $(O_FILES)

include $(GOROOT)/src/Make.cmd

lexer.$O:
	$(GC) -o lexer.$O src/lexer/lexer.go src/lexer/token.go

parser.$O:
	$(GC) -o parser.$O src/parser/parser.go

