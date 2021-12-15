package tests

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"sed/interpreter"
	"sed/parser"
	"strings"
	"testing"
)

func TestRunning_1(t *testing.T) {
	ctx := parser.ParseExpression("/<p>/!b;:x N;/<\\/p>/!bx;s/\n//g;", false)

	assert := assert.New(t)
	ast := parser.BeginLabel.Next().Next()
	l, ok := ast.(*parser.Label)
	assert.True(ok)
	assert.True(l.Name == "x")

	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug: true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	assert.True(AssertContent("a\n<p>bc</p>\ne\n<p>fg</p>\nh"))

}

func TestRunning_2_1(t *testing.T) {
	ctx := parser.ParseExpression("1h;1!H;x;l;/<\\/p>/{s/\n//gp;n};/<p>/!{p;n};x;", false)

	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		InhibitPrint: true,
		Debug: true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a$\na\n<p>\\nb$\n<p>\\nb\\nc$\n<p>\\nb\\nc\\n</p>$\n<p>bc</p>\ne\n<p>\\nf$\n<p>\\nf\\ng$\n<p>\\nf\\ng\\n</p>$\n<p>fg</p>\nh"))
}

func TestRunning_2_2(t *testing.T) {
	ctx := parser.ParseExpression("H;x;l;/<\\/p>/{s/\n//gp;n};/<p>/!{p;n};x;", false)

	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		InhibitPrint: true,
		Debug: true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a$\na\n<p>\\nb$\n<p>\\nb\\nc$\n<p>\\nb\\nc\\n</p>$\n<p>bc</p>\ne\n<p>\\nf$\n<p>\\nf\\ng$\n<p>\\nf\\ng\\n</p>$\n<p>fg</p>\nh"))
}

func TestRunning_3(t *testing.T) {
	ctx := parser.ParseExpression(":x; s_(<p>.*)\n(.*</p>)_\\1\\2_;tx; /<p>[^\n]*<\\/p>/n; N; bx", false)

	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
		Debug: true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("a\n<p>\nb\nc\n</p>\ne\n<p>\nf\ng\n</p>\nh")))
	(assert.New(t)).True(AssertContent("a\n<p>bc</p>\ne\n<p>fg</p>\nh"))
}

func TestG(t *testing.T) {
	ctx := parser.ParseExpression(":x;N;/(.*)\"(.*)\n *(.*)\"(.*)/!bx;s//\\1\n\"\\2 \\3\"\n\\4/", false)

	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		InhibitPrint: false,
		Debug: true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("And he said \"This too\n   shall pass\" end")))
	(assert.New(t)).True(AssertContent("And he said \n\"This too shall pass\"\n end"))
}

func TestG2(t *testing.T) {
	ctx := parser.ParseExpression(":x;$!{N;bx}; s/\n//g", false)

	ec := &parser.ExecutionContext{
		CurrentAST:   ctx.AST,
		InhibitPrint: false,
		Debug: true,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("And he said \"This too\n   shall pass\" end")))
	(assert.New(t)).True(AssertContent("And he said \"This too   shall pass\" end"))

}