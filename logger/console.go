package logger

import (
	"fmt"
)

type consoleLogger struct {
	name string
}

func NewConsoleLogger() *consoleLogger {

	return &consoleLogger{}
}

func (c *consoleLogger) SetName(name string) {
	c.name = name
}

func (c *consoleLogger) Log(msg string) {
	if c.name != "" {
		msg = fmt.Sprintf("[%s] %s", c.name, msg)
	}
	fmt.Println(msg)
}

func (c *consoleLogger) Logf(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	c.Log(msg)
}
