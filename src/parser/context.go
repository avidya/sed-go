package parser

import (
	"fmt"
	"os"
)

type ExecutionContext struct {
	HoldSpace      string
	PatternSpace   string
	CurrentLineNum int
	RangeCount     int
	LastSubResult  bool
	InRange        bool
	HitEnd         bool
	InhibitPrint   bool
	LastPattern    string
	Empty          bool
	CurrentAST     ASTNode
	StreamMode     bool
	Debug          bool
	DebugFile      *os.File
}

func PrintContent(content string, ec *ExecutionContext) {
	fmt.Print(content)
	if ec.Debug {
		if ec.DebugFile == nil {
			ec.DebugFile, _ = os.Create("/tmp/sed-go-output")
		}
		ec.DebugFile.WriteString(content)
	}
}
