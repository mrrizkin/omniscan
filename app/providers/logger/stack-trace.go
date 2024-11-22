package logger

import (
	"fmt"
	"runtime/debug"
	"strings"
)

type StackFrame struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

func stackTrace() []StackFrame {
	stack := debug.Stack()
	lines := strings.Split(string(stack), "\n")
	frames := make([]StackFrame, 0)

	// Skip first 9 lines to remove this function and the logging calls
	for i := 9; i < len(lines)-1; i += 2 {
		// Skip empty lines
		if len(lines[i]) == 0 {
			continue
		}

		// Parse function name
		funcName := strings.TrimSpace(lines[i])

		// Parse file and line number
		if i+1 >= len(lines) {
			break
		}
		fileLine := strings.TrimSpace(lines[i+1])

		var frame StackFrame
		frame.Function = funcName

		// Parse file:line format
		if strings.Contains(fileLine, ":") {
			parts := strings.Split(fileLine, ":")
			frame.File = parts[0]
			if len(parts) > 1 {
				fmt.Sscanf(parts[1], "%d", &frame.Line)
			}
		}

		frames = append(frames, frame)
	}

	return frames
}
