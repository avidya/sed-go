package tests

import (
	"github.com/stretchr/testify/assert"
	"sed/parser"
	"testing"
)

func TestParseExpression_Empty(t *testing.T) {

	ctx := parser.ParseExpression("", false)
	assert := assert.New(t)
	assert.True(ctx.AST == parser.BeginLabel)
	assert.True(ctx.AST.Next() == parser.EndLabel)
	assert.True(ctx.AST.Next().Next() == nil)
}

func TestParseExpression_Substitute_SuccessWoutFlags(t *testing.T) {

	ctx := parser.ParseExpression("s/a/b/", false)
	assert := assert.New(t)
	assert.True(ctx.AST == parser.BeginLabel)
	stat, ok := ctx.AST.Next().(*parser.Statement)
	assert.True(ok)
	assert.True(len(stat.Addresses) == 0)
	assert.True(len(stat.Functions) == 1)
	sub, ok := stat.Functions[0].(*parser.Substitute)
	assert.True(ok)
	//assert.True(sub.Pattern != nil && sub.Pattern.String() == "a")
	assert.True(sub.Replacement == "b")
	assert.True(!sub.Global)
	assert.True(!sub.Print)
	assert.True(sub.Occurrence == -1)
}

func TestParseExpression_Substitute_SuccessWithFlags(t *testing.T) {
	ctx := parser.ParseExpression("s/a/b/pg2", false)
	sub := ctx.AST.Next().(*parser.Statement).Functions[0].(*parser.Substitute)
	assert := assert.New(t)

	assert.True(sub.Global)
	assert.True(sub.Print)
	assert.True(sub.Occurrence == 2)
}

func TestParseExpression_Substitute_SuccessWithFlags2(t *testing.T) {
	ctx := parser.ParseExpression("s/a/b/p2g", false)
	sub := ctx.AST.Next().(*parser.Statement).Functions[0].(*parser.Substitute)
	assert := assert.New(t)

	assert.True(sub.Global)
	assert.True(sub.Print)
	assert.True(sub.Occurrence == 2)
}