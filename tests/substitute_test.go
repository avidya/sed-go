package tests

import (
	"bufio"
	"github.com/avidya/sed-go/interpreter"
	"github.com/avidya/sed-go/parser"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestEvalSub_Simple(t *testing.T) {
	ctx := parser.ParseExpression("s/a/b/", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("abc")))
	assert := assert.New(t)
	assert.True(ec.PatternSpace == "bbc")
}

func TestEvalSub_Global(t *testing.T) {
	ctx := parser.ParseExpression("s/a/b/g", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("aaa")))
	assert := assert.New(t)
	assert.True(ec.PatternSpace == "bbb")
}

func TestEvalSub_Occurrence(t *testing.T) {
	ctx := parser.ParseExpression("s/a/b/2", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("aaa")))
	assert := assert.New(t)
	assert.True(ec.PatternSpace == "aba")

	ctx = parser.ParseExpression("s/a/b/2g", false)
	ec = &parser.ExecutionContext{
		CurrentAST: ctx.AST,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("aaa")))
	assert.True(ec.PatternSpace == "abb")
}

func TestEvalSub_BackreferenceAndAmpersand(t *testing.T) {
	ctx := parser.ParseExpression("s/(ab).?/\\1_&z/", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("abc-vabz")))
	assert := assert.New(t)
	assert.True(ec.PatternSpace == "ab_abcz-vabz")

	// escape the Ampersand
	ctx = parser.ParseExpression("s/(ab).?/\\1_\\&z/", false)
	ec = &parser.ExecutionContext{
		CurrentAST: ctx.AST,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("abc-vabz")))
	assert.True(ec.PatternSpace == "ab_&z-vabz")

	// escape the Ampersand, and replace the second match part
	ctx = parser.ParseExpression("s/(ab).?/\\1_\\&z/2", false)
	ec = &parser.ExecutionContext{
		CurrentAST: ctx.AST,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("abc-vabz")))
	assert.True(ec.PatternSpace == "abc-vab_&z")

	// complex global
	ctx = parser.ParseExpression("s/(ab).?/\\1_&z/g", false)
	ec = &parser.ExecutionContext{
		CurrentAST: ctx.AST,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("abc-vabz")))
	assert.True(ec.PatternSpace == "ab_abcz-vab_abzz")

	ctx = parser.ParseExpression("s/(ab).?/\\\\1_&z/", false)
	ec = &parser.ExecutionContext{
		CurrentAST: ctx.AST,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("abc-vabz")))
	assert.True(ec.PatternSpace == "\\1_abcz-vabz")

}

func TestEvalSub_DelimiterOtherThanSlash(t *testing.T) {
	ctx := parser.ParseExpression("s*(ab).?*\\1_&z*", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("abc-vabz")))
	assert := assert.New(t)
	assert.True(ec.PatternSpace == "ab_abcz-vabz")
}

func TestEvalSub_IgnoreCase(t *testing.T) {
	ctx := parser.ParseExpression("s/B/BEE/i", false)
	ec := &parser.ExecutionContext{
		CurrentAST: ctx.AST,
	}
	interpreter.Eval(ec, *bufio.NewScanner(strings.NewReader("abc")))
	assert := assert.New(t)
	assert.True(ec.PatternSpace == "aBEEc")
}
