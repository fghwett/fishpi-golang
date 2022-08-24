package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"fishpi/logger"
)

type Handler struct {
	oldTopic  string        // 旧标题
	red       *WsMsgReply   // 拼手气红包、平分红包
	gesture   *WsMsgReply   // 猜拳红包
	heartbeat *WsMsgReply   // 心跳红包
	own       *WsMsgReply   // 专属红包
	lastest   *WsMsgReply   // 最近一条消息
	cache     []*WsMsgReply // 消息缓存

	cacheNum int
	sdk      *Sdk
	logger   logger.Logger
}

func NewHandler(cacheNum int, sdk *Sdk, logger logger.Logger) *Handler {
	h := &Handler{
		cacheNum: cacheNum,
		sdk:      sdk,
		logger:   logger,
	}

	return h
}

func (h *Handler) HandleMsg(data interface{}) {
	bytes, ok := data.([]byte)
	if !ok {
		h.logger.Logf("data is not []byte: %v", string(bytes))
		return
	}

	msg := &WsMsgReply{}
	if err := json.Unmarshal(bytes, &msg); err != nil {
		log.Printf("parse message error: %s, body: %s\n", err, string(bytes))
		return
	}
	msg.Parse()
	h.filterMessage(msg)

	content := msg.Msg()
	if msg.Type == WsMsgTypeOnline {
		if content == h.oldTopic {
			return
		}
		h.oldTopic = content
	} else if msg.Type == WsMsgTypeRevoke {
		content = fmt.Sprintf("有人撤回了一条消息 消息内容不知道 %s %s", msg.OId, msg.UserAvatarURL210)
		for _, v := range h.cache {
			if msg.OId == v.OId {
				content = fmt.Sprintf("有人撤回了一条消息：%s", v.Msg())
				break
			}
		}
	}

	h.logger.Log(content)
}

func (h *Handler) filterMessage(msg *WsMsgReply) {
	if msg.IsRedPacketMsg() {
		switch msg.RedPackageInfo.Type {
		case RedPacketTypeSpecify:
			h.own = msg
		case RedPacketTypeRockPaperScissors:
			h.gesture = msg
		case RedPacketTypeHeartbeat:
			h.heartbeat = msg
		default:
			h.red = msg
		}
		return
	}

	if msg.Type == WsMsgTypeMsg {
		h.cache = append(h.cache, msg)
		if len(h.cache) >= h.cacheNum {
			removeNum := len(h.cache) - h.cacheNum
			h.cache = h.cache[removeNum:]
		}

		if msg.UserName == h.sdk.username {
			h.lastest = msg
		}
	}
}

func (h *Handler) HandleWsStatusMsg(data interface{}) {
	str, ok := data.(string)
	if !ok {
		h.logger.Logf("data is not string: %v", str)
		return
	}
	h.logger.Log(str)
}

func (h *Handler) KeepLive() <-chan string {
	c := make(chan string)
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case _ = <-ticker.C:
				c <- "-hb-"
			}
		}
	}()

	return c
}

func (h *Handler) Watch() {
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
		go h.handleCommand(recv)
	}
}

func (h *Handler) handleCommand(cmd string) {
	if cmd == "0" || cmd == "1" || cmd == "2" || cmd == "3" { // 抢红包 0-普通红包(拼手气 平分) 1-3猜拳红包 4-心跳红包 5-专属红包
		h.handleReceiveRedPacket(cmd)
	} else if cmd == "revoke" { // 撤回最近一条消息
		h.handleRevokeLastMessage()
	} else if cmd == "repeat" { // 重复最近一条消息
		h.handleRepeatLastMessage()
	} else {
		h.logger.Logf("无效指令：%s", cmd)
	}
}

func (h *Handler) handleReceiveRedPacket(gesture string) {
	var red *WsMsgReply
	switch gesture {
	case "0":
		red = h.red
	case "1", "2", "3":
		red = h.gesture
	case "4":
		red = h.heartbeat
	case "5":
		red = h.own
	}
	if red == nil {
		return
	}

	result, err := h.sdk.OpenRedPacket(red.OId, gesture)
	if err != nil {
		h.logger.Logf("打开红包%s失败 %s", red.OId, err)
		return
	}
	h.logger.Log(result)
}

func (h *Handler) handleRevokeLastMessage() {
	if h.lastest == nil {
		h.logger.Log("您最近还没有讲话")
		return
	}
	if err := h.sdk.RevokeMsg(h.lastest.OId); err != nil {
		h.logger.Log(err.Error())
		return
	}
	h.logger.Log("撤回消息操作成功")
}

func (h *Handler) handleRepeatLastMessage() {
	if h.cache == nil || len(h.cache) == 0 {
		return
	}
	msg := h.cache[len(h.cache)-1]
	if err := h.sdk.SendMsg(msg.Md); err != nil {
		h.logger.Log(err.Error())
	}
}
