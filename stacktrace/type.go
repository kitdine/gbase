package stacktrace

import (
	"fmt"
	"runtime"
	"strings"
)

// CurrentStack 获取当前 goroutine 的结构化堆栈（跳过 skip 帧）main 中 默认为2
func CurrentStack(skip int) GoroutineStackInfo {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(skip+2, pcs) // +2: skip CurrentStack 和 Callers 本身
	frames := runtime.CallersFrames(pcs[:n])

	var trace []string
	for {
		frame, more := frames.Next()
		trace = append(trace, frame.Function)
		trace = append(trace, fmt.Sprintf("\t%s:%d", frame.File, frame.Line))
		if !more {
			break
		}
	}

	return GoroutineStackInfo{
		GoroutineName: "goroutine current [running]:", // 无法准确获取 ID，可自定义说明
		Traces:        trace,
	}
}

// 获取所有 goroutine 的结构化堆栈（并跳过当前 goroutine 顶部 N 帧） main中 默认为 3
func AllGoroutineStacks(skipCurrentTopFrames int) []GoroutineStackInfo {
	buf := make([]byte, 1<<20)
	n := runtime.Stack(buf, true)
	raw := string(buf[:n])

	parts := strings.Split(raw, "\n\n")
	stacks := make([]GoroutineStackInfo, 0, len(parts))

	for i, part := range parts {
		lines := strings.Split(part, "\n")
		if len(lines) < 2 {
			continue
		}

		// 如果是当前 goroutine（即第一个），可跳过顶部帧
		if i == 0 && len(lines) > 2+skipCurrentTopFrames*2 {
			lines = append(lines[:1], lines[1+skipCurrentTopFrames*2:]...)
		}

		stacks = append(stacks, GoroutineStackInfo{
			GoroutineName: lines[0],
			Traces:        lines[1:],
		})
	}

	return stacks
}
