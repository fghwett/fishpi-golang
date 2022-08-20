package eventHandler

import (
	"testing"
	"time"

	"fishpi/logger"
)

func TestEventHandler(t *testing.T) {
	l := logger.NewConsoleLogger()
	eh := NewEventHandler("test", l)

	var (
		func1 = func(data interface{}) {
			t.Log("func1 receive: ", data)
		}
		func2 = func(data interface{}) {
			t.Log("func2 receive: ", data)
		}
		func3 = func(data interface{}) {
			t.Log("func3 receive: ", data)
		}
	)

	const eventLog EventType = "event_log"

	eh.Sub(eventLog, func1)
	eh.Sub(eventLog, func2)
	eh.Sub(eventLog, func3)

	eh.Pub(eventLog, "hello world")

	time.Sleep(time.Second)
}
