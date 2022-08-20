package eventHandler

type EventType string

const (
	WsConnected       = "ws-connected"
	WsClosed          = "ws-closed"
	WsReconnectedFail = "ws-reconnected-fail"
	WsMsg             = "ws-msg"
)

type EventHandler interface {
	Pub(eventType EventType, data interface{})
	Sub(eventType EventType, callback func(data interface{}))
}
