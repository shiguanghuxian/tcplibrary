/*
 * @Author: 时光弧线
 * @Date: 2017-12-30 11:55:33
 * @Last Modified by: 时光弧线
 * @Last Modified time: 2017-12-30 13:02:01
 */
package tcplibrary

import (
	"fmt"
	"log"
)

/* tcp 内部日志打印 */

// Logger 日志记录
type Logger struct {
	debug bool
}

var globalLogger *Logger

func init() {
	globalLogger = new(Logger)
	globalLogger.debug = true
}

// SetDebug 设置是否为debug模式
func (l *Logger) SetDebug(debug bool) {
	l.debug = debug
}

// Infof 打印错误 Info
func (l *Logger) Infof(format string, a ...interface{}) {
	if l.debug == false {
		return
	}
	format = fmt.Sprintf("INFO: %s\n", format)
	log.Printf(format, a...)
}

// Warnf 打印错误 Warn
func (l *Logger) Warnf(format string, a ...interface{}) {
	if l.debug == false {
		return
	}
	format = fmt.Sprintf("WARN: %s\n", format)
	log.Printf(format, a...)
}

// Errorf 打印错误 Error
func (l *Logger) Errorf(format string, a ...interface{}) {
	format = fmt.Sprintf("ERROR: %s\n", format)
	log.Printf(format, a...)
}

// Fatalf 打印错误 Fatal
func (l *Logger) Fatalf(format string, a ...interface{}) {
	format = fmt.Sprintf("FATAL: %s\n", format)
	log.Printf(format, a...)
}
