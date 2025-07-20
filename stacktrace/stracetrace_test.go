package stacktrace

import (
	"github.com/kitdine/gbase/log"
	"go.uber.org/zap"
	"testing"
)

func TestAllGoroutineStacks(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			log.GetLogger().Fatal("panic occurred",
				zap.Any("recover", r),
				zap.Any("current_stack", CurrentStack(2)),
				zap.Any("all_goroutines", AllGoroutineStacks(3)), // 👈 跳过顶部封装帧
			)
		}
	}()
	panic("something bad happened")
}
