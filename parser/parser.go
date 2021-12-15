package parser

/**
The total grammar about sed is summarized as follows(in EBNF format):

grammar sed;

statements
  : statement ( ';' statement )* ;

statement
  : ':' label
  | address?  function | '{' functions '}'
;

functions
  : function ( ';'? function )* ;

address
  : firstAddress ( ',' secondAddress )? '!'? ;

firstAddress
  : absoluteAddress ;

absoluteAddress
  : [1-9][0-9]*
  | '$'
  | pattern
;

secondAddress
  : absoluteAddress |  relativeAddress ;

relativeAddress
  : '+' [1-9][0-9]* ;

pattern
  : '/' .+ '/'
  | '\' (delimiter) .+ \1
;

delimiter
  : [^\\n] ;

function
  : 'b' label?
  | 't' label?
  | 's' (delimiter) .* \1 .* \1 flags*
  | [nNgGhHxlpD]
;

label
  : [a-z]+ ；

flags
  : [[0-9]gp]

 Author: kozz.gaof
 Date: 2021/11/22
*/
import (
	"fmt"
	"regexp"
	"unicode"
)

const RegExpFlagsMatchNewline = "(?s)"
const RegExpFlagsCaseInsensitive = "(?i)"

type AnalysisContext struct {
	Source         []rune
	position       int
	current        rune
	AST            ASTNode
	LastPatternStr string
	StreamMod      bool
}

type ASTNode interface {
	Next() ASTNode
	SetNext(ast ASTNode)
}

type Label struct {
	Name     string
	NextNode *ASTNode
}

func (l *Label) Next() ASTNode {
	if l.NextNode != nil {
		return *l.NextNode
	} else {
		return nil
	}
}

func (l *Label) SetNext(ast ASTNode) {
	l.NextNode = &ast
}

type Statement struct {
	Addresses []AddrPattern
	Functions []Function
	NextNode  *ASTNode
}

func (s *Statement) Next() ASTNode {
	return *s.NextNode
}

func (s *Statement) SetNext(ast ASTNode) {
	s.NextNode = &ast
}

var BeginLabel = &Label{Name: "BEGIN"}
var EndLabel = &Label{Name: "END"}

func ParseExpression(source string, streamMod bool) *AnalysisContext {
	ctx := &AnalysisContext{
		Source:    []rune(source),
		position:  -1,
		AST:       BeginLabel,
		StreamMod: streamMod,
	}
	ctx.statements()
	ctx.AST.SetNext(EndLabel)
	ctx.AST = BeginLabel
	return ctx
}

func (ctx *AnalysisContext) statements() {
	for token := ctx.NextToken(); token != EOF; token = ctx.NextToken() {
		ctx.statement()
		if ctx.NextToken() == SEMICOLON {
			ctx.match(SEMICOLON)
		}
	}
}

func (ctx *AnalysisContext) statement() {
	token := ctx.NextToken()
	var ast ASTNode
	if token == COLON {
		ctx.match(COLON)
		ast = ctx.label()
		ctx.AST.SetNext(ast)
		ctx.AST = ast
	} else {
		ast = &Statement{}
		ctx.AST.SetNext(ast)
		ctx.AST = ast
		ctx.addresses()
		ctx.functions()
	}
}

func (ctx *AnalysisContext) label() ASTNode {
	l := ctx.ConsecutiveChars(';', ' ')
	if len(l) == 0 {
		panic("lacks a label")
	} else {
		return &Label{
			Name: l,
		}
	}
}

func (ctx *AnalysisContext) addresses() {
	ctx.firstAddress()
	if token := ctx.NextToken(); token == COMMA {
		ctx.match(COMMA)
		ctx.secondAddress()
	}
	if ctx.NextToken() == EXCLAMATION {
		ctx.match(EXCLAMATION)
		for _, addr := range ctx.AST.(*Statement).Addresses {
			addr.SetNegative(true)
		}
	}
}

func (ctx *AnalysisContext) firstAddress() {
	ctx.absoluteAddr(true)
}

func (ctx *AnalysisContext) secondAddress() {
	token := ctx.NextToken()
	if token = ctx.NextToken(); token == PLUS {
		ctx.relativeAddr()
	} else {
		ctx.absoluteAddr(false)
	}
}

func (ctx *AnalysisContext) absoluteAddr(first bool) {
	if token := ctx.NextToken(); unicode.IsDigit(token) {
		ctx.digits()
	} else if token == DOLLAR {
		ctx.match(DOLLAR)
		ctx.ends()
	} else if token == SLASH {
		ctx.match(SLASH)
		ctx.pattern('/', first)
		ctx.match(SLASH)
	} else if token == BACK_SLASH {
		ctx.match(BACK_SLASH)
		delimiter := ctx.NextToken()
		ctx.match(delimiter)
		ctx.pattern(delimiter, first)
		ctx.match(delimiter)
	}
	// no default situation should be processed, since address is optional.
}

func (ctx *AnalysisContext) digits() {
	lineNum := ConsecutiveInt(ctx)
	s := ctx.AST.(*Statement)
	s.Addresses = append(s.Addresses, &LinePattern{
		LineNumber: lineNum,
	})
}

func (ctx *AnalysisContext) ends() {
	s := ctx.AST.(*Statement)
	s.Addresses = append(s.Addresses, &EndPattern{})
}

func (ctx *AnalysisContext) pattern(delimiter rune, first bool) {
	if l := patternStr(ctx, delimiter); len(l) > 0 {
		ctx.LastPatternStr = l
		s := ctx.AST.(*Statement)
		if first {
			s.Addresses = append(s.Addresses, &BeginPattern{PatternStr: l})
		} else {
			s.Addresses = append(s.Addresses, &StopPattern{PatternStr: l})
		}
	}
}

func patternStr(ctx *AnalysisContext, delimiter rune) string {
	str := ctx.ConsecutiveChars(delimiter)
	// find escaping sign
	for (len(str) > 1 && str[len(str)-1] == '\\' && str[len(str)-2] != '\\') || (len(str) == 1 && str[0] == '\\') {
		b := []rune(str)
		str = string(append(b[:len(b)-1], ctx.NextToken()))
		ctx.match(ctx.NextToken())
		str += patternStr(ctx, delimiter) //  recursive finding
	}
	//fmt.Println(str)
	return str
}

func (ctx *AnalysisContext) relativeAddr() {
	ctx.match(PLUS)
	count := ConsecutiveInt(ctx)
	s := ctx.AST.(*Statement)
	s.Addresses = append(s.Addresses, &RelativePattern{
		Count: count,
	})
}

func (ctx *AnalysisContext) functions() {
	if token := ctx.NextToken(); token == L_BRACE {
		ctx.match(L_BRACE)
		for ctx.NextToken() != R_BRACE {
			ctx.function()
			if ctx.NextToken() == SEMICOLON {
				ctx.match(SEMICOLON)
			}
		}
		ctx.match(R_BRACE)
	} else {
		ctx.function()
	}
}

func (ctx *AnalysisContext) function() {
	s := ctx.AST.(*Statement)
	switch ctx.NextToken() {
	case 'N':
		ctx.match('N')
		s.Functions = append(s.Functions, &NextAppend{})
	case 'n':
		ctx.match('n')
		s.Functions = append(s.Functions, &Next{})
	case 'G':
		ctx.match('G')
		s.Functions = append(s.Functions, &HoldSpaceAppendToPatternSpace{})
	case 'g':
		ctx.match('g')
		s.Functions = append(s.Functions, &HoldSpaceToPatternSpace{})
	case 'H':
		ctx.match('H')
		s.Functions = append(s.Functions, &PatternAppendToHoldSpace{})
	case 'h':
		ctx.match('h')
		s.Functions = append(s.Functions, &PatternSpaceToHoldSpace{})
	case 'x':
		ctx.match('x')
		s.Functions = append(s.Functions, &Exchange{})
	case 's':
		ctx.match('s')
		delimiter := ctx.NextToken()
		ctx.match(delimiter)
		sub := &Substitute{Occurrence: -1}
		str := patternStr(ctx, delimiter)
		if len(str) == 0 {
			if len(ctx.LastPatternStr) == 0 {
				panic("no previous regular expression")
			} else {
				str = ctx.LastPatternStr
			}
		}
		ctx.match(delimiter)
		sub.Replacement = patternStr(ctx, delimiter)
		ctx.match(delimiter)
		s.Functions = append(s.Functions, sub)
		ctx.flags()
		if sub.CaseInsensitive {
			sub.Pattern = regexp.MustCompile(RegExpFlagsCaseInsensitive + RegExpFlagsMatchNewline + str)
		} else {
			sub.Pattern = regexp.MustCompile(RegExpFlagsMatchNewline + str)
		}
	case 'D':
		ctx.match('D')
		s.Functions = append(s.Functions, &DeleteToNewline{})
	case 'd':
		ctx.match('d')
		s.Functions = append(s.Functions, &Delete{})
	case 'l':
		ctx.match('l')
		s.Functions = append(s.Functions, &Debug{})
	case 'p':
		ctx.match('p')
		s.Functions = append(s.Functions, &Print{})
	case 'b':
		ctx.match('b')
		label := ctx.ConsecutiveChars(';', ' ', '}')
		s.Functions = append(s.Functions, &Jump{Name: label})
	case 't':
		ctx.match('t')
		label := ctx.ConsecutiveChars(';', ' ')
		s.Functions = append(s.Functions, &ConditionalJump{Name: label})
	}
}

func (ctx *AnalysisContext) flags() {
	for ctx.NextToken() != SEMICOLON && ctx.NextToken() != EOF && ctx.NextToken() != R_BRACE {
		functions := ctx.AST.(*Statement).Functions
		sub := functions[len(functions)-1].(*Substitute)
		if token := ctx.NextToken(); unicode.IsDigit(token) {
			if sub.Occurrence != -1 {
				panic("multiple number options to `s' command")
			} else {
				sub.Occurrence = ConsecutiveInt(ctx)
			}
		} else if token == 'g' {
			if sub.Global {
				panic(" multiple `g' options to `s' command")
			} else {
				sub.Global = true
			}
			ctx.match('g')
		} else if token == 'p' {
			if sub.Print {
				panic(" multiple `p' options to `s' command")
			} else {
				sub.Print = true
			}
			ctx.match('p')
		} else if token == 'i' {
			sub.CaseInsensitive = true
			ctx.match('i')
		} else {
			panic("unknown option to `s'")
		}
	}
}

// match : everytime this func matches a given char, it'll call Consume() func subsequently, to move the underlying cursor one char forward.
func (ctx *AnalysisContext) match(char rune) {
	if ctx.NextToken() == char {
		ctx.Consume()
	} else {
		panic(fmt.Sprintf("expecting %v; found %v\n", char, ctx.NextToken()))
	}
}
