// +build js

package frugal

type safeLogger struct {
}

func logger() *safeLogger {
	return &safeLogger{}
}

func (log *safeLogger) Error(msg string) {
	println("error", msg)
}

func (log *safeLogger) Debug(msg string) {
	println("debug", msg)
}

func (log *safeLogger) Warn(msg string) {
	println("warn", msg)
}

func (log *safeLogger) Print(msg string) {
	println("print", msg)
}

func (log *safeLogger) Info(msg string) {
	println("info", msg)
}
