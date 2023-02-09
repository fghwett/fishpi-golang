package eventHandler

type EventType string

const (
	WsConnected       = "ws-connected"
	WsClosed          = "ws-closed"
	WsReconnectedFail = "ws-reconnected-fail"
	WsMsg             = "ws-msg"
	WsSend            = "ws-send"

	ElvesStick = `elves-stick` // 召唤小飞棍
)

type EventHandler interface {
	Pub(eventType EventType, data interface{})
	Sub(eventType EventType, callback func(data interface{}))
}
