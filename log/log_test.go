package log

import (
	"testing"
)

func TestLogger(t *testing.T) {
	base := ConfigBase{
		MaxSize:        500,  // 但文件大小，单位：MB
		MaxBackups:     100,  // 备份文件个数
		MaxAge:         14,   // 最大保存天数
		Compress:       true, // 备份文件是否压缩保存
		JSONFormat:     true, // 日志打印格式，是否启用json 格式
		ShowLineNumber: true, // 是否显示打印位置信息，类、行号等
	}
	config := GlobalConfig{
		Level:       "info",          // 日志级别
		LogFileName: "logs/root.log", // 日志文件 全路径名
		ConfigBase:  base,
	}
	kafka := ChildConfig{
		LoggerName:  "kafka",
		Level:       "info",
		LogFileName: "logs/kafka.log",
		ConfigBase:  base,
	}
	InitLogger(config, kafka)

	GetLogger().Info("this is global logger")
	GetLoggerWithName("kafka").Info("this is kafka logger")
}
