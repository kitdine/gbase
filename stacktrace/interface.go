package stacktrace

type GoroutineStackInfo struct {
	GoroutineName string   `json:"goroutine,omitempty"`
	Traces        []string `json:"traces,omitempty"`
}
