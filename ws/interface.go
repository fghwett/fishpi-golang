package ws

type Websocket interface {
	Send(msg interface{})
	Start() error
	Stop() error
}
