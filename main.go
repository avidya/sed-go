package main

import (
	"bufio"
	"flag"
	"os"
	"sed/interpreter"
	"sed/parser"
)

func main() {
	inhibitPrint := flag.Bool("n", false, "inhibit printing")
	streamMod := flag.Bool("s", false, "work in pseudo-stream mode")
	flag.Parse()
	expression := flag.Args()

	ec := &parser.ExecutionContext{
		CurrentAST:   parser.ParseExpression(expression[0], *streamMod).AST,
		InhibitPrint: *inhibitPrint,
		StreamMode:   *streamMod,
	}

	if *streamMod {
		interpreter.EvalStream(ec, *bufio.NewScanner(bufio.NewReader(os.Stdin)))
	} else {
		interpreter.Eval(ec, *bufio.NewScanner(bufio.NewReader(os.Stdin)))
	}
}
