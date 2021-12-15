package tests

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"sed/interpreter"
	"sed/parser"
	"strings"
	"testing"
)

func TestStreamMod(t *testing.T) {
	ctx := parser.ParseExpression("/<p>/,\\x</p>x s_\n__g;", true)

	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		InhibitPrint: false,
		StreamMode: true,
		Debug: true,
	}
	interpreter.EvalStream(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a\n<p>bc</p>\ne\n<p>fg</p>\nh"))
}

func TestStreamMod2_1(t *testing.T) {
	ctx := parser.ParseExpression("1,$ s/\n//g", true)

	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		StreamMode: true,
		InhibitPrint: false,
		Debug: true,
	}
	interpreter.EvalStream(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a<p>bc</p>e<p>fg</p>h"))
}

func TestStreamMod2_2(t *testing.T) {
	ctx := parser.ParseExpression("s/\n//g", true)

	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		StreamMode: true,
		InhibitPrint: false,
		Debug: true,
	}
	interpreter.EvalStream(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a<p>bc</p>e<p>fg</p>h"))
}

func TestStreamMod3(t *testing.T) {
	ctx := parser.ParseExpression("2,8 s/\n//g", true)

	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		StreamMode: true,
		InhibitPrint: false,
		Debug: true,
	}
	interpreter.EvalStream(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a\n<p>bc</p>e<p>f\ng\n</p>\nh"))
}

func TestStreamMod4(t *testing.T) {
	ctx := parser.ParseExpression("2,+6 s/\n//g", true)

	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		StreamMode: true,
		InhibitPrint: false,
		Debug: true,
	}
	interpreter.EvalStream(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a\n<p>bc</p>e<p>f\ng\n</p>\nh"))
}

