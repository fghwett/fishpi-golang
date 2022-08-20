package ws

type Websocket interface {
	Send(msg []byte)
	Start() error
	Stop() error
}
