package tests

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"sed/interpreter"
	"sed/parser"
	"strings"
	"testing"
)

func TestSingleAddress(t *testing.T) {
	ctx := parser.ParseExpression("2 s/a/b/", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug:      true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\na\na")))
	(assert.New(t)).True(AssertContent("a\nb\na"))
}

func TestSingleAddressNegative(t *testing.T) {
	ctx := parser.ParseExpression("2! s/a/b/", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug:      true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\na\na")))
	(assert.New(t)).True(AssertContent("b\na\nb"))
}

func TestRangeMatch(t *testing.T) {
	ctx := parser.ParseExpression("\\x<p>x,\\y</p>yp", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug:      true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a\n<p>\n<p>\nb\nb\nc\nc\n</p>\n</p>\ne\n<p>\n<p>\nf\nf\ng\ng\n</p>\n</p>\nh"))
}

func TestRangeMatchSecondMismatch(t *testing.T) {
	ctx := parser.ParseExpression("\\x<p>x,/<\\/P>/p", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug:      true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a\n<p>\n<p>\nb\nb\nc\nc\n</p>\n</p>\ne\ne\n<p>\n<p>\nf\nf\ng\ng\n</p>\n</p>\nh\nh"))
}

func TestRangeMatchSecondMismatchNegative(t *testing.T) {
	ctx := parser.ParseExpression("\\x<p>x,/<\\/P>/!p", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug:      true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a\na\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh"))
}

func TestRangeMSingle(t *testing.T) {
	ctx := parser.ParseExpression("\\xbxs//B/p", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug:      true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a\n<p>\nB\nB\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh"))
}

func TestRangeMSingle2(t *testing.T) {
	ctx := parser.ParseExpression("\\b\\bb s//B/p", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug:      true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a\n<p>\nB\nB\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh"))
}

func TestRangeMSingleError(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			assert.New(t).True(true)
		}
	}()
	parser.ParseExpression("\\b\\bbs//B/p'", false)
	assert.New(t).True(false)
}

func TestRelativeRange(t *testing.T) {
	ctx := parser.ParseExpression("2,+1 s/.*/XX/", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug:      true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\nb\nc\nd")))
	(assert.New(t)).True(AssertContent("a\nXX\nXX\nd"))

}

func TestRelativeRangeNegative(t *testing.T) {
	ctx := parser.ParseExpression("2,+1! s/.*/XX/", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug:      true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\nb\nc")))
	(assert.New(t)).True(AssertContent("XX\nb\nc"))
}

func TestEndRange(t *testing.T) {
	ctx := parser.ParseExpression("$ s/.*/XX/", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug:      true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\nb\nc")))
	(assert.New(t)).True(AssertContent("a\nb\nXX"))
}

func TestEndRangeNegative(t *testing.T) {
	ctx := parser.ParseExpression("$! s/.*/XX/", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug:      true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("XX\nXX\nc")))
}
