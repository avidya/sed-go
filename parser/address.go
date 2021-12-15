package parser

import (
	"fmt"
	"regexp"
)

type AddrPattern interface {
	Match(ec *ExecutionContext) bool
	SetNegative(n bool)
	Negative() bool
}

type BeginPattern struct {
	PatternStr string
	pattern    *regexp.Regexp
	n          bool
}

func (bp *BeginPattern) Match(ec *ExecutionContext) bool {
	if bp.pattern == nil {
		bp.pattern = regexp.MustCompile(fmt.Sprintf("%s(.*)(%s.*)", RegExpFlagsMatchNewline, bp.PatternStr))
	}
	result := bp.pattern.FindAllStringSubmatch(ec.PatternSpace, -1)
	if len(result) > 0 {
		if ec.StreamMode {
			if !ec.InhibitPrint {
				PrintContent(result[0][1], ec) // shouldn't print newline feed here
			}
			ec.PatternSpace = result[0][2]
		}
		return true
	} else {
		return false
	}
}

func (bp *BeginPattern) Negative() bool {
	return bp.n
}

func (bp *BeginPattern) SetNegative(n bool) {
	bp.n = n
}

type StopPattern struct {
	PatternStr string
	pattern    *regexp.Regexp
	n          bool
}

func (sp *StopPattern) Match(ec *ExecutionContext) bool {
	if sp.pattern == nil {
		sp.pattern = regexp.MustCompile(fmt.Sprintf("%s(.*%s)(.*)", RegExpFlagsMatchNewline, sp.PatternStr))
	}
	result := sp.pattern.FindAllStringSubmatch(ec.PatternSpace, -1)

	if len(result) > 0 {
		if ec.StreamMode {
			ec.PatternSpace = result[0][1]
			ec.HoldSpace = result[0][2]
		}
		return true
	} else {
		return false
	}
}

func (sp *StopPattern) Negative() bool {
	return sp.n
}

func (sp *StopPattern) SetNegative(n bool) {
	sp.n = n
}

type LinePattern struct {
	LineNumber int
	n          bool
}

func (lp *LinePattern) Match(ec *ExecutionContext) bool {
	return lp.LineNumber == ec.CurrentLineNum
}

func (lp *LinePattern) Negative() bool {
	return lp.n
}

func (lp *LinePattern) SetNegative(n bool) { lp.n = n }

type RelativePattern struct {
	Count int
	n     bool
}

func (rp *RelativePattern) Match(ec *ExecutionContext) bool {
	return rp.Count == ec.RangeCount
}

func (rp *RelativePattern) Negative() bool {
	return rp.n
}

func (rp *RelativePattern) SetNegative(n bool) { rp.n = n }

type EndPattern struct {
	n bool
}

func (ep *EndPattern) Match(ec *ExecutionContext) bool {
	return ec.HitEnd
}

func (ep *EndPattern) Negative() bool {
	return ep.n
}

func (ep *EndPattern) SetNegative(n bool) { ep.n = n }
