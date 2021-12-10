package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

const (
	Normal = iota
	Continue
	Break
)

type Function interface {
	Call(ec *ExecutionContext, scanner *bufio.Scanner) int
}

type NextAppend struct {
}

func (n *NextAppend) Call(ec *ExecutionContext, scanner *bufio.Scanner) int {
	if !ec.HitEnd {
		if len(ec.PatternSpace) == 0 {
			ec.PatternSpace = scanner.Text()
		} else {
			ec.PatternSpace = ec.PatternSpace + "\n" + scanner.Text()
		}
		ec.Empty = false
		ec.HitEnd = !scanner.Scan()
		ec.CurrentLineNum++
		return Normal
	} else {
		return Break
	}
}

type Next struct {
}

func (n2 *Next) Call(ec *ExecutionContext, scanner *bufio.Scanner) int {
	if !ec.HitEnd {
		if !ec.InhibitPrint {
			PrintContent(ec.PatternSpace+"\n", ec)
		}
		ec.PatternSpace = scanner.Text()
		ec.Empty = false
		ec.HitEnd = !scanner.Scan()
		ec.CurrentLineNum++
		return Normal
	} else {
		return Break
	}
}

type Jump struct {
	Name string
}

func (b *Jump) Call(ec *ExecutionContext, _ *bufio.Scanner) int {
	if len(b.Name) == 0 {
		ec.CurrentAST = EndLabel
		return Continue
	} else {
		var ast ASTNode
		for ast = BeginLabel; ast != EndLabel; ast = ast.Next() {
			if label, ok := ast.(*Label); ok && label.Name == b.Name {
				ec.CurrentAST = ast
				return Continue
			}
		}
		panic(fmt.Sprintf("can't find label for jump to `%s'", b.Name))
	}
}

type ConditionalJump struct {
	Name string
}

func (t *ConditionalJump) Call(ec *ExecutionContext, scanner *bufio.Scanner) int {
	if !ec.LastSubResult {
		return Normal
	} else {
		return (&Jump{t.Name}).Call(ec, scanner)
	}
}

type Debug struct {
}

func (l *Debug) Call(ec *ExecutionContext, _ *bufio.Scanner) int {
	str := strings.ReplaceAll(ec.PatternSpace, "\\", "\\\\")
	str = strings.ReplaceAll(str, "\n", "\\n")
	str = strings.ReplaceAll(str, "\a", "\\a")
	str = strings.ReplaceAll(str, "\f", "\\f")
	str = strings.ReplaceAll(str, "\r", "\\r")
	str = strings.ReplaceAll(str, "\t", "\\t")
	str = strings.ReplaceAll(str, "\v", "\\v")
	PrintContent(str+"$\n", ec)
	return Normal
}

type Print struct {
}

func (p *Print) Call(ec *ExecutionContext, _ *bufio.Scanner) int {
	PrintContent(ec.PatternSpace+"\n", ec)
	return Normal
}

type Substitute struct {
	Pattern         *regexp.Regexp
	Replacement     string
	Occurrence      int
	Global          bool
	Print           bool
	CaseInsensitive bool
}

func (s *Substitute) Call(ec *ExecutionContext, _ *bufio.Scanner) int {
	result := s.Pattern.FindAllStringSubmatchIndex(ec.PatternSpace, -1)
	ec.LastSubResult = false
	if len(result) > 0 {
		psBytes := []byte(ec.PatternSpace)
		psBuf := bytes.NewBuffer([]byte{})
		lastIndex := 0
		for index, group := range result {
			if index+1 >= s.Occurrence {
				psBuf.Write(psBytes[lastIndex:group[0]])
				lastIndex = group[1]
				replBuf := bytes.NewBuffer([]byte{})
				var escape = false
				for _, b := range []byte(s.Replacement) {
					if escape {
						if b > 48 && b < 58 { // is a digit
							i := (b - 48) << 1
							replBuf.Write(psBytes[group[i]:group[i+1]])
						} else if b == '\\' || b == '&' {
							replBuf.WriteByte(b)
						} else if b == 'n' {
							replBuf.WriteByte('\n')
						}
						escape = false
					} else if b == '\\' {
						escape = true
					} else if b == '&' {
						replBuf.Write(psBytes[group[0]:group[1]])
					} else {
						replBuf.WriteByte(b)
					}
				}
				psBuf.Write(replBuf.Bytes())
				ec.LastSubResult = true
				if !s.Global {
					break
				}
			}
		}
		psBuf.Write(psBytes[lastIndex:])
		ec.PatternSpace = psBuf.String()
		if s.Print {
			PrintContent(ec.PatternSpace+"\n", ec)
		}
	}
	return Normal
}

type Exchange struct {
}

func (x *Exchange) Call(ec *ExecutionContext, _ *bufio.Scanner) int {
	ec.HoldSpace, ec.PatternSpace = ec.PatternSpace, ec.HoldSpace
	return Normal
}

type PatternAppendToHoldSpace struct {
}

func (h *PatternAppendToHoldSpace) Call(ec *ExecutionContext, _ *bufio.Scanner) int {
	if len(ec.HoldSpace) == 0 {
		// I hate to append a `\n' even for the first line! To avoid this, one usually has to write something like: 1h;1!H; ...
		ec.HoldSpace = ec.PatternSpace
	} else {
		ec.HoldSpace = ec.HoldSpace + "\n" + ec.PatternSpace
	}
	return Normal
}

type PatternSpaceToHoldSpace struct {
}

func (h2 *PatternSpaceToHoldSpace) Call(ec *ExecutionContext, _ *bufio.Scanner) int {
	ec.HoldSpace = ec.PatternSpace
	return Normal
}

type HoldSpaceAppendToPatternSpace struct {
}

func (g *HoldSpaceAppendToPatternSpace) Call(ec *ExecutionContext, _ *bufio.Scanner) int {
	ec.PatternSpace = ec.PatternSpace + "\n" + ec.HoldSpace
	return Normal
}

type HoldSpaceToPatternSpace struct {
}

func (g2 *HoldSpaceToPatternSpace) Call(ec *ExecutionContext, _ *bufio.Scanner) int {
	ec.PatternSpace = ec.HoldSpace
	return Normal
}

type Delete struct {
}

func (d *Delete) Call(ec *ExecutionContext, _ *bufio.Scanner) int {
	ec.PatternSpace = ""
	ec.Empty = true
	return Break
}

type DeleteToNewline struct {
}

func (d2 *DeleteToNewline) Call(ec *ExecutionContext, _ *bufio.Scanner) int {
	fragments := strings.Split(ec.PatternSpace, "\n")
	ec.PatternSpace = strings.Join(fragments[1:], "\n")
	return Break
}
