package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Interface 日志接口。
type Interface interface {
	Debug(message interface{}, fields ...interface{})
	Info(message string, fields ...interface{})
	Warn(message string, fields ...interface{})
	Error(message interface{}, fields ...interface{})
	Fatal(message interface{}, fields ...interface{})
}

// Level 日志等级。
type Level int

const (
	ErrorLevel Level = iota
	WarnLevel
	InfoLevel
	DebugLevel
)

// Logger 基于标准库 log 的日志实现。
type Logger struct {
	logger *log.Logger
	level  Level
}

var _ Interface = (*Logger)(nil)

// New 创建 Logger。
func New(level string) *Logger {
	var l Level

	switch strings.ToLower(level) {
	case "error":
		l = ErrorLevel
	case "warn":
		l = WarnLevel
	case "info":
		l = InfoLevel
	case "debug":
		l = DebugLevel
	default:
		l = InfoLevel
	}

	logger := log.New(os.Stdout, "", log.LstdFlags)

	return &Logger{
		logger: logger,
		level:  l,
	}
}

// Debug 记录 Debug 日志。
func (l *Logger) Debug(message interface{}, fields ...interface{}) {
	if !l.shouldLog(DebugLevel) {
		return
	}

	l.msg("debug", message, fields...)
}

// Info 记录 Info 日志。
func (l *Logger) Info(message string, fields ...interface{}) {
	if !l.shouldLog(InfoLevel) {
		return
	}

	l.log(message, fields...)
}

// Warn 记录 Warn 日志。
func (l *Logger) Warn(message string, fields ...interface{}) {
	if !l.shouldLog(WarnLevel) {
		return
	}

	l.log(message, fields...)
}

// Error 记录 Error 日志。
func (l *Logger) Error(message interface{}, fields ...interface{}) {
	if l.level == DebugLevel {
		l.Debug(message, fields...)
	}

	if !l.shouldLog(ErrorLevel) {
		return
	}

	l.msg("error", message, fields...)
}

// Fatal 记录 Fatal 日志并退出进程。
func (l *Logger) Fatal(message interface{}, fields ...interface{}) {
	if !l.shouldLog(ErrorLevel) {
		os.Exit(1)
	}

	l.msg("fatal", message, fields...)

	os.Exit(1)
}

func (l *Logger) log(message string, fields ...interface{}) {
	formatted := message
	if len(fields) != 0 {
		formatted = fmt.Sprintf(message, fields...)
	}

	if formatted == "" {
		return
	}

	l.logger.Println(formatted)
}

func (l *Logger) msg(level string, message interface{}, fields ...interface{}) {
	switch msg := message.(type) {
	case error:
		l.log(fmt.Sprintf("[%s] %s", level, msg.Error()), fields...)
	case string:
		l.log(fmt.Sprintf("[%s] %s", level, msg), fields...)
	default:
		l.log(fmt.Sprintf("[%s] message %v has unknown type %T", level, message, msg), fields...)
	}
}

func (l *Logger) shouldLog(target Level) bool {
	return l.level >= target
}
