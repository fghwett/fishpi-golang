package ws

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"fishpi/eventHandler"
	"fishpi/logger"
)

type ws struct {
	addr              string
	reconnectInterval int
	breakReconnect    bool

	client *websocket.Conn

	sendChan chan []byte
	readChan chan []byte

	event  eventHandler.EventHandler
	logger logger.Logger

	ctx    context.Context
	cancel context.CancelFunc
}

func NewWs(addr string, reconnectInterval int, event eventHandler.EventHandler, logger logger.Logger) *ws {
	w := &ws{
		addr:              addr,
		reconnectInterval: reconnectInterval,

		sendChan: make(chan []byte, 1024),
		readChan: make(chan []byte, 1024),

		event:  event,
		logger: logger,
	}

	go w.handle()

	return w
}

func (w *ws) Start() error {
	if w.client != nil {
		return errors.New("旧连接尚未断开")
	}
	w.breakReconnect = true
	return w.conn()
}

func (w *ws) Stop() error {
	if w.client == nil {
		return errors.New("没有检测到连接")
	}
	w.breakReconnect = false
	return w.client.Close()
}

func (w *ws) conn() error {

	c, _, err := websocket.DefaultDialer.Dial(w.addr, nil)
	if err != nil {
		return err
	}
	w.event.Pub(eventHandler.WsConnected, fmt.Sprintf("Websocket Connect Success"))

	w.client = c
	w.client.SetPongHandler(func(appData string) error {
		log.Printf("recevice pong, %s\n", appData)
		return nil
	})

	w.client.SetCloseHandler(func(code int, text string) error {
		log.Printf("addr: %s\ncode: %d\ntext: %s\n", w.addr, code, text)
		w.event.Pub(eventHandler.WsClosed, fmt.Sprintf("Websocket closed: \ncode: %d\ntext: %s\naddr: %s", code, text, w.addr))
		w.reconnect()

		return nil
	})

	w.ctx, w.cancel = context.WithCancel(context.Background())
	go w.read()
	go w.write()

	return nil
}

func (w *ws) reConn() {
	time.Sleep(time.Duration(w.reconnectInterval) * time.Second)

	if err := w.conn(); err != nil {
		w.logger.Logf("conn %s error: %s", w.addr, err)
		w.event.Pub(eventHandler.WsReconnectedFail, fmt.Sprintf("Websocket Reconnected failed\nerror: %s\naddr: %s", err, w.addr))
		go w.reConn()
	}
}

func (w *ws) Send(msg []byte) {
	w.sendChan <- msg
}

func (w *ws) write() {
	for {
		select {
		case msg := <-w.sendChan:
			if err := w.client.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("%s write message error: %s", w.addr, err)
				return
			}
		case <-w.ctx.Done():
			log.Printf("stop write progress\n")
			return
		}
	}
}

func (w *ws) read() {
	for {
		_, message, err := w.client.ReadMessage()
		if err != nil {
			log.Printf("%s read message error: %s, stop read message", w.addr, err)
			w.reconnect()
			return
		}
		w.readChan <- message
	}
}

func (w *ws) handle() {
	for {
		select {
		case msg := <-w.readChan:
			//log.Printf("receive message: %s\n", string(msg))
			w.event.Pub(eventHandler.WsMsg, msg)
		}
	}
}

func (w *ws) reconnect() {
	w.cancel()
	if err := w.client.Close(); err != nil {
		log.Printf("close ws connect error: %s\n", err)
	}
	w.client = nil

	if w.breakReconnect {
		go w.reConn()
	}
}
