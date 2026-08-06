// Package logger 提供按月滚动的结构化日志功能。
// debug 模式输出文本格式到终端，生产模式同时输出 JSON 到终端和日志文件。
//
// 使用示例：
//
//	logger.InitLog(logger.Config{
//	    DebugMode: os.Getenv("MODE") == "debug",
//	    LogDir:    "logs",
//	})
//	logger.Info("服务启动", "port", 9000)
//	logger.Error("数据库错误", "err", err)
package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

// Config 日志初始化配置
type Config struct {
	// DebugMode 为 true 时只输出到终端（文本格式，含 Debug 级别）
	// 为 false 时同时输出到终端和日志文件（JSON 格式，Info 级别及以上）
	DebugMode bool

	// LogDir 日志文件存放目录，默认 "logs"
	// 仅在 DebugMode=false 时有效
	LogDir string
}

var (
	instance    *slog.Logger
	mu          sync.Mutex
	currentFile *os.File
	cfg         Config
	initialized bool
)

// InitLog 初始化日志，应在程序启动时调用一次（main 或 init 中）
func InitLog(c Config) {
	mu.Lock()
	defer mu.Unlock()

	if c.LogDir == "" {
		c.LogDir = "logs"
	}
	cfg = c

	if cfg.DebugMode {
		instance = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	} else {
		if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
			panic(fmt.Sprintf("logger: 创建日志目录失败: %v", err))
		}
		instance = newProductionLogger()
		scheduleRotate() // 启动定时切换
	}

	slog.SetDefault(instance)
	initialized = true
}

// newProductionLogger 创建生产模式 logger（同时写终端和文件）
func newProductionLogger() *slog.Logger {
	month := time.Now().Format("2006-01")
	f := openLogFile(month)
	multi := io.MultiWriter(os.Stdout, f)
	return slog.New(slog.NewJSONHandler(multi, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// openLogFile 打开指定月份的日志文件，关闭旧文件
func openLogFile(month string) *os.File {
	filename := fmt.Sprintf("%s/app-%s.log", cfg.LogDir, month)
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(fmt.Sprintf("logger: 打开日志文件失败: %v", err))
	}
	if currentFile != nil {
		currentFile.Close()
	}
	currentFile = f
	return f
}

// scheduleRotate 计算距离下个月1号00:00的时间差，到点切换日志文件后再次调度
func scheduleRotate() {
	now := time.Now()
	// 下个月1号 00:00:00
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	duration := nextMonth.Sub(now)

	time.AfterFunc(duration, func() {
		mu.Lock()
		instance = newProductionLogger()
		slog.SetDefault(instance)
		mu.Unlock()

		scheduleRotate() // 递归调度下一次
	})
}

// get 返回当前 logger 实例，如果未初始化则 panic
func get() *slog.Logger {
	if !initialized {
		panic("logger: 未调用 InitLog() 初始化")
	}
	return instance
}

// Info 记录 Info 级别日志
func Info(msg string, args ...any) {
	get().Info(msg, args...)
}

// Error 记录 Error 级别日志
func Error(msg string, args ...any) {
	get().Error(msg, args...)
}

// Warn 记录 Warn 级别日志
func Warn(msg string, args ...any) {
	get().Warn(msg, args...)
}

// Debug 记录 Debug 级别日志（仅 DebugMode=true 时输出）
func Debug(msg string, args ...any) {
	get().Debug(msg, args...)
}

// With 返回带有固定字段的子 logger，适合在特定模块里使用
// 例如：reqLogger := logger.With("module", "invoice")
func With(args ...any) *slog.Logger {
	return get().With(args...)
}
