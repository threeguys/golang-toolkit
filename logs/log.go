//
//  Copyright 2020 Ray Cole
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//  http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.
//
package logs

import (
    "fmt"
    "runtime/debug"
)

type LogLevels string

const (
    NONE = "NONE"
    DEBUG = "DEBUG"
    INFO = "INFO"
    WARN = "WARN"
    ERROR = "ERROR"
)

type Logger interface {
    Log(level LogLevels, msg string, args ... interface{})
    Info(msg string, args ... interface{})
    Warn(msg string, args ... interface{})
    Debug(msg string, args ... interface{})
    Error(msg string, args ... interface{})
}

type SimpleLogger struct {
    name string
    parent Logger
}

func (l *SimpleLogger) Log(level LogLevels, msg string, args ... interface{}) {
    if l.parent != nil {
        l.parent.Log(level, msg, args...)
    } else {
        fmt.Printf("%s: %s\n", level, fmt.Sprintf(msg, args...))
    }
}

func (l *SimpleLogger) Info(msg string, args ... interface{}) {
    l.Log(INFO, msg, args...)
}

func (l *SimpleLogger) Warn(msg string, args ... interface{}) {
    l.Log(WARN, msg, args...)
}

func (l *SimpleLogger) Debug(msg string, args ... interface{}) {
    l.Log(DEBUG, msg, args...)
}

func (l *SimpleLogger) Error(msg string, args ... interface{}) {
    l.Log(ERROR, msg, args...)
}

var rootLogger = Logger(&SimpleLogger{ name:   "(root)" })
func GetLogger(name string) Logger {
    return &SimpleLogger{
        name:   name,
        parent: rootLogger,
    }
}

func SetRootLogger(log Logger) {
    rootLogger = log
}

type WriterLogger struct {
    Writer func([]byte) error
}

func (w WriterLogger) Log(level LogLevels, msg string, args ...interface{}) {
    output := fmt.Sprintf(msg, args...)
    err := w.Writer([]byte(fmt.Sprintf("[%s] - %s\n", level, output)))
    if err != nil {
        fmt.Errorf("Unable to write to logger: %s", string(debug.Stack()))
    }
}

func (w WriterLogger) Info(msg string, args ...interface{}) {
    w.Log(INFO, msg, args...)
}

func (w WriterLogger) Warn(msg string, args ...interface{}) {
    w.Log(WARN, msg, args...)
}

func (w WriterLogger) Debug(msg string, args ...interface{}) {
    w.Log(DEBUG, msg, args...)
}

func (w WriterLogger) Error(msg string, args ...interface{}) {
    w.Log(ERROR, msg, args...)
}


