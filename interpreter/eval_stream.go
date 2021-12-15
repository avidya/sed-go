package interpreter

import (
	"bufio"
	p "sed/parser"
)

func EvalStream(ec *p.ExecutionContext, scanner bufio.Scanner) {
	for ec.HitEnd = !scanner.Scan(); !ec.HitEnd; {
		if ec.InRange {
			(&p.NextAppend{}).Call(ec, &scanner)
			ec.RangeCount++
		} else {
			ec.PatternSpace = scanner.Text()
			ec.Empty = false
			ec.HitEnd = !scanner.Scan()
			ec.CurrentLineNum++
		}
		if stat, ok := p.BeginLabel.Next().(*p.Statement); !ok {
			panic("label is not allowed in stream mode")
		} else {
			for matchStream(ec, stat.Addresses) && len(ec.PatternSpace) > 0 && (!ec.InRange || matchStream(ec, stat.Addresses)) {
				for _, function := range stat.Functions {
					function.Call(ec, &scanner)
				}
				if !ec.InhibitPrint && !ec.InRange {
					p.PrintContent(ec.PatternSpace+"\n", ec)
				}
				ec.PatternSpace = ec.HoldSpace
			}
		}
		if !ec.InhibitPrint && !ec.InRange && len(ec.PatternSpace) > 0 {
			p.PrintContent(ec.PatternSpace+"\n", ec)
		}
	}
}

func matchStream(ec *p.ExecutionContext, addresses []p.AddrPattern) bool {
	switch len(addresses) {
	case 0:
		ec.InRange = true
		return ec.CurrentLineNum == 1 || ec.HitEnd
	case 1:
		panic("single address does not work in stream mode")
	case 2:
		if !ec.InRange {
			ec.InRange = addresses[0].Match(ec)
			return ec.InRange
		} else {
			ec.InRange = !addresses[1].Match(ec)
			return !ec.InRange
		}
	default:
		panic("only 0~2 addresses are allowed to be specified")
	}
}
