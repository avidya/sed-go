package tests

import (
	"bufio"
	"github.com/avidya/sed-go/interpreter"
	"github.com/avidya/sed-go/parser"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestFunc(t *testing.T) {
	ctx := parser.ParseExpression("hgNhnxGs/\n//g;s/abc/ABC/p", false)

	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		InhibitPrint: false,
		Debug:        true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\nb\nc")))
	(assert.New(t)).True(AssertContent("a\nb\nABC\nABC"))
}

func TestDeleteToNewline(t *testing.T) {
	ctx := parser.ParseExpression("NNlDl", false)
	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		InhibitPrint: false,
		Debug:        true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\nb\nc")))
	(assert.New(t)).True(AssertContent("a\\nb\\nc$\nb\nc"))
}

func TestDelete(t *testing.T) {
	ctx := parser.ParseExpression("NNldl", false)
	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		InhibitPrint: false,
		Debug:        true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\nb\nc\nd")))
	(assert.New(t)).True(AssertContent("a\\nb\\nc$\nd"))
}

func TestLineDelete_1(t *testing.T) {
	ctx := parser.ParseExpression("3d", false)

	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug:      true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a\n<p>\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh"))
}

func TestLineDelete_2(t *testing.T) {
	ctx := parser.ParseExpression("3,$d", false)

	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug:      true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a\n<p>"))
}
