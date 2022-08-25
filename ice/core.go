package ice

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"fishpi/logger"
)

type core struct {
	ck       string
	username string
	uid      string

	ch chan []byte

	logger logger.Logger
}

func NewCore(ck, username, uid string, logger logger.Logger) *core {
	c := &core{
		ck:       ck,
		username: username,
		uid:      uid,

		ch: make(chan []byte, 1000),

		logger: logger,
	}

	return c
}

func (c *core) login() {
	body, _ := json.Marshal(&ExchangeMsg{
		Type: TypeSetUser,
		User: c.username,
		Ck:   c.ck,
		Uid:  c.uid,
	})

	c.ch <- body
}

func (c *core) HandleMsg(data interface{}) {
	bytes, ok := data.([]byte)
	if !ok {
		c.logger.Logf("data is not []byte: %v", string(bytes))
		return
	}

	msg := &ExchangeMsg{}
	if err := json.Unmarshal(bytes, &msg); err != nil {
		log.Printf("parse message error: %s, body: %s\n", err, string(bytes))
		return
	}

	m := msg.Msg
	m = strings.ReplaceAll(m, "<br>", "\n")
	m = strings.ReplaceAll(m, "<summary>", "")
	m = strings.ReplaceAll(m, "</summary>", "")
	m = strings.ReplaceAll(m, "<details>", "\n")
	m = strings.ReplaceAll(m, "</details>", "")

	if msg.Type == TypeAll {
		c.logger.Log(m)
		c.login()
	} else if msg.Type == TypeGameMsg {
		if msg.VipLv != 0 {
			c.logger.Logf("%s %s", msg.Level(), m)
		} else {
			c.logger.Log(m)
		}
	} else {
		c.logger.Log(string(bytes))
	}
}

func (c *core) HandleWsStatusMsg(data interface{}) {
	str, ok := data.(string)
	if !ok {
		c.logger.Logf("data is not string: %v", str)
		return
	}
	c.logger.Log(str)
}

func (c *core) Watch() {
	for {
		var buf [1024]byte
		read := bufio.NewReader(os.Stdin)
		m, err := read.Read(buf[:])
		if err != nil {
			log.Printf("read error: %s\n", err)
			continue
		}
		recv := strings.Split(string(buf[:m]), "\n")[0]
		recv = strings.TrimSpace(recv)
		if recv == "" {
			continue
		}
		go c.handleCommand(recv)
	}
}

func (c *core) handleCommand(cmd string) {
	body, _ := json.Marshal(&ExchangeMsg{
		Type: TypeGameMsg,
		Ck:   c.ck,
		Msg:  cmd,
	})
	c.ch <- body
}

func (c *core) KeepLive() <-chan []byte {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case _ = <-ticker.C:
				body, _ := json.Marshal(&ExchangeMsg{Type: TypeHb})
				c.ch <- body
			}
		}
	}()

	return c.ch
}
