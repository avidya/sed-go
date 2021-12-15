package interpreter

import (
	"bufio"
	p "sed/parser"
)

func Eval(ec *p.ExecutionContext, scanner bufio.Scanner) {
	for ec.HitEnd = !scanner.Scan(); !ec.HitEnd; {

		ec.PatternSpace = scanner.Text()
		ec.Empty = false
		ec.HitEnd = !scanner.Scan()
		ec.CurrentLineNum++

		for ec.CurrentAST = p.BeginLabel; ec.CurrentAST != p.EndLabel; {
			if stat, ok := ec.CurrentAST.(*p.Statement); ok && matchAddress(ec, ec.CurrentAST.(*p.Statement).Addresses) {
				status := p.Normal
				for _, function := range stat.Functions {
					if status = function.Call(ec, &scanner); status == p.Continue || status == p.Break {
						break
					}
				}
				if status == p.Break {
					break
				} else if status == p.Continue {
					continue
				}
			}
			ec.CurrentAST = ec.CurrentAST.Next()
		}
		if !ec.InhibitPrint && !ec.Empty {
			p.PrintContent(ec.PatternSpace+"\n", ec)
		}
	}
}

func matchAddress(ec *p.ExecutionContext, addresses []p.AddrPattern) bool {
	switch len(addresses) {
	case 0:
		return true
	case 1:
		return xor(addresses[0].Match(ec), addresses[0].Negative())
	case 2:
		if !ec.InRange {
			ec.InRange = addresses[0].Match(ec)
			return xor(ec.InRange, addresses[0].Negative())
		} else {
			r := ec.InRange // include
			ec.RangeCount++
			ec.InRange = !addresses[1].Match(ec)
			return xor(r, addresses[1].Negative())
		}
	default:
		panic("only 0~2 addresses are allowed to be specified")
	}
}

func xor(b1 bool, b2 bool) bool {
	return b1 != b2
}
