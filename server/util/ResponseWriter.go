package util

import (
	"fmt"
	"io"
	"strings"
)

// ResponseWriter is an interface to handling output to either a human or a computer using a fixed
// API
type ResponseWriter interface {
	Println(level int, args ...interface{})
	Printf(level int, format string, args ...interface{})
}

const (
	ERR_INFORMATIONAL   = 100
	ERR_OK              = 200
	ERR_NOT_FOUND       = 404
	ERROR               = 500
	ERR_UNKNOWN_COMMAND = 501
)

const (
	KEY_WRITER = "responseWriter"
)

// HumanResponseWriter is a ResponseWriter for human consumption
type HumanResponseWriter struct {
	out io.Writer
}

func NewHumanResponseWriter(out io.Writer) ResponseWriter {
	return &HumanResponseWriter{out: out}
}

func (w *HumanResponseWriter) write(s string) {
	s = strings.ReplaceAll(s, "\n", "\r\n")
	_, _ = w.out.Write([]byte(s[:]))
}

func (w *HumanResponseWriter) Println(_ int, args ...interface{}) {
	w.write(fmt.Sprint(args...))
}

func (w *HumanResponseWriter) Printf(_ int, format string, args ...interface{}) {
	w.write(fmt.Sprintf(format, args...))
}

// APIResponseWriter is a ResponseWriter for api consumption
type APIResponseWriter struct {
	out io.Writer
}

func NewAPIResponseWriter(out io.Writer) ResponseWriter {
	return &APIResponseWriter{out: out}
}

func WriteString(level int, s string, f func(s string)) {
	prefix := fmt.Sprintf("%03d", level)
	a := strings.Split(s, "\n")
	l := len(a) - 1
	for i, e := range a {
		sep := "-"
		if i >= l {
			sep = " "
		}
		f(prefix + sep + e)
	}
}

func (w *APIResponseWriter) Println(level int, args ...interface{}) {
	WriteString(level, fmt.Sprint(args...), func(s string) {
		_, _ = w.out.Write([]byte(s))
	})
}

func (w *APIResponseWriter) Printf(level int, format string, args ...interface{}) {
	WriteString(level, fmt.Sprintf(format, args...), func(s string) {
		_, _ = w.out.Write([]byte(s))
	})
}
