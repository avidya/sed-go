package parser

import (
	"strconv"
	"unicode"
)

//NextToken : whitespace will be automatically skipped
func (ctx *AnalysisContext) NextToken() rune {
	if ctx.position == -1 {
		ctx.Consume()
	}
	switch ctx.current {
	case EOF:
		return EOF
	case ' ':
		ctx.Consume()
		return ctx.NextToken()
	default:
		return ctx.current
	}
}

func (ctx *AnalysisContext) Consume() {
	ctx.position++
	if ctx.position >= len(ctx.Source) {
		ctx.current = EOF
	} else {
		ctx.current = ctx.Source[ctx.position]
	}
}

func ConsecutiveInt(ctx *AnalysisContext) int {
	var l []rune
	for ; unicode.IsDigit(ctx.current); ctx.Consume() {
		l = append(l, ctx.current)
	}
	i, _ := strconv.Atoi(string(l))
	return i
}

// ConsecutiveChars : whitespace won't be skipped
func (ctx *AnalysisContext) ConsecutiveChars(chars ...rune) string {
	var l []rune
	for ; !contains(chars, ctx.current) && ctx.current != EOF; ctx.Consume() {
		l = append(l, ctx.current)
	}
	return string(l)
}

func contains(chars []rune, char rune) bool {
	for _, c := range chars {
		if c == char {
			return true
		}
	}
	return false
}
