# go base tools

## Logger

Example
```go
import (
    "github.com/kitdine/gbase/log"
)

config := log.Config{
Level:          "info", // 日志级别
LogFileName:    "/path/to/root.log", // 日志文件 全路径名
MaxSize:        500, // 但文件大小，单位：MB
MaxBackups:     100, // 备份文件个数
MaxAge:         14, // 最大保存天数
Compress:       true, // 备份文件是否压缩保存
JSONFormat:     true, // 日志打印格式，是否启用json 格式
ShowLineNumber: true, // 是否显示打印位置信息，类、行号等
}
err := log.InitLogger(config)
if err != nil {
fmt.Println(err.Error())
}

// 使用logger
logger.GetLogger().Info("Hello World！！！")
logger.GetLogger().Error("Hello World！！！")
```
