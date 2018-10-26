/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:55:33
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2018-10-21 14:41:00
 */
package tcplibrary

import (
	"fmt"
	"log"
)

/* tcp 内部日志打印 */

// TCPLibraryLogger 日志接口，可以使用第三方日志库，符合接口即可
type TCPLibraryLogger interface {
	Infof(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Fatalf(format string, a ...interface{})
}

// tcplibrary 全局日志对象
var globalLogger TCPLibraryLogger

// SetGlobalLogger 设置日志对象
func SetGlobalLogger(logger TCPLibraryLogger) {
	if logger != nil {
		globalLogger = logger
	} else {
		globalLogger.Errorf("设置tcplibrary日志对象不能是nil")
	}
}

func init() {
	defaultLogger := new(Logger)
	defaultLogger.debug = true
	globalLogger = defaultLogger
}

// Logger 日志记录
type Logger struct {
	debug bool
}

// SetDefaultDebug 设置是否为debug模式 - 只对默认日志库生效
func (l *Logger) SetDefaultDebug(debug bool) {
	l.debug = debug
}

// Infof 打印错误 Info
func (l *Logger) Infof(format string, a ...interface{}) {
	if l.debug == false {
		return
	}
	format = fmt.Sprintf("INFO: %s\n", format)
	log.Output(2, fmt.Sprintf(format, a...))
}

// Warnf 打印错误 Warn
func (l *Logger) Warnf(format string, a ...interface{}) {
	if l.debug == false {
		return
	}
	format = fmt.Sprintf("WARN: %s\n", format)
	log.Output(2, fmt.Sprintf(format, a...))
}

// Errorf 打印错误 Error
func (l *Logger) Errorf(format string, a ...interface{}) {
	format = fmt.Sprintf("ERROR: %s\n", format)
	log.Output(2, fmt.Sprintf(format, a...))
}

// Fatalf 打印错误 Fatal
func (l *Logger) Fatalf(format string, a ...interface{}) {
	format = fmt.Sprintf("FATAL: %s\n", format)
	log.Output(2, fmt.Sprintf(format, a...))
}
