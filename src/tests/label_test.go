package tests

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"sed/interpreter"
	"sed/parser"
	"strings"
	"testing"
)

func TestConditionalJump(t *testing.T) {
	ctx := parser.ParseExpression(":x;s/a/A/;tx", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("aaa")))
	assert.New(t).True(ec.PatternSpace == "AAA")
}
