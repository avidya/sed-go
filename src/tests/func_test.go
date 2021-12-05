package tests

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"sed/interpreter"
	"sed/parser"
	"strings"
	"testing"
)

func TestFunc(t *testing.T) {
	ctx := parser.ParseExpression("hgNhnxGs/\n//g;s/abc/ABC/p", false)

	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		InhibitPrint: false,
		Debug: true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\nb\nc")))
	(assert.New(t)).True(AssertContent("a\nb\nABC\nABC"))
}

func TestDeleteToNewline(t *testing.T) {
	ctx := parser.ParseExpression("NNlDl", false)
	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		InhibitPrint: false,
		Debug: true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\nb\nc")))
	(assert.New(t)).True(AssertContent("a\\nb\\nc$\nb\nc"))
}

func TestDelete(t *testing.T) {
	ctx := parser.ParseExpression("NNldl", false)
	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		InhibitPrint: false,
		Debug: true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\nb\nc\nd")))
	(assert.New(t)).True(AssertContent("a\\nb\\nc$\nd"))
}