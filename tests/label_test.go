package tests

import (
	"bufio"
	"github.com/avidya/sed-go/interpreter"
	"github.com/avidya/sed-go/parser"
	"github.com/stretchr/testify/assert"
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
