package core

import (
	"encoding/json"
	"fishpi/eventHandler"
	"time"
)

type Core struct {
	oldTopic  *WsMsgReply   // 旧标题
	red       *WsMsgReply   // 拼手气红包、平分红包
	gesture   *WsMsgReply   // 猜拳红包
	heartbeat *WsMsgReply   // 心跳红包
	own       *WsMsgReply   // 专属红包
	lastest   *WsMsgReply   // 最近一条消息
	cache     []*WsMsgReply // 消息缓存

	msgChannel   chan *WsMsgReply
	showMsgCache []*WsMsgReply

	cacheNum int
	token    string
	sdk      *Sdk
	eh       eventHandler.EventHandler
}

func NewCore(cacheNum int, token string, sdk *Sdk, eh eventHandler.EventHandler) *Core {
	c := &Core{
		cacheNum: cacheNum,
		token:    token,
		sdk:      sdk,
		eh:       eh,
	}

	c.init()
	c.KeepLive()
	return c
}

func (c *Core) init() {
	//data, err := c.sdk.ChatRecordPage(1)
	//if err != nil {
	//	//c.showMsg(fmt.Sprintf("获取历史聊天记录失败 %s", err.Error()))
	//	return
	//}
	//sort.Slice(data, func(i, j int) bool {
	//	t1, _ := time.Parse("2006-01-02 15:04:05", data[i].Time)
	//	t2, _ := time.Parse("2006-01-02 15:04:05", data[j].Time)
	//	return t1.Before(t2)
	//})
	//for _, v := range data {
	//	//content := strings.TrimPrefix(strings.TrimSuffix(v.Content, "</p>"), "<p>")
	//	//c.showMsg(fmt.Sprintf("%s %s(%s): %s", v.Time[11:], v.UserNickname, v.UserName, content))
	//}
}

// SendPublicMsg 发送消息
func (c *Core) SendPublicMsg(content string) error {
	return c.sdk.SendMsg(content)
}

// GetUserInfo 获取用户信息
func (c *Core) GetUserInfo(username string) string {
	return c.sdk.UserInfo(username)
}

// OpenRedPacket 打开红包
func (c *Core) OpenRedPacket(oId, gesture string) (string, error) {
	return c.sdk.OpenRedPacket(oId, gesture)
}

func (c *Core) HandleMsg(data interface{}) {
	bytes, ok := data.([]byte)
	if !ok {
		return
	}

	msg := &WsMsgReply{}
	if err := json.Unmarshal(bytes, &msg); err != nil {
		return
	}
	msg.Parse()
	c.filterMessage(msg)

	//content := msg.Msg()
	//if msg.Type == WsMsgTypeOnline {
	//	if c.oldTopic == nil {
	//		c.oldTopic = msg
	//		return
	//	}
	//	if msg.Msg() == c.oldTopic.Msg() {
	//		return
	//	}
	//	c.oldTopic = msg
	//} else if msg.Type == WsMsgTypeRevoke {
	//	content = fmt.Sprintf("有人撤回了一条消息 消息内容不知道 %s %s", msg.OId, msg.UserAvatarURL210)
	//	for _, v := range c.cache {
	//		if msg.OId == v.OId {
	//			content = fmt.Sprintf("有人撤回了一条消息：%s", v.Msg())
	//			break
	//		}
	//	}
	//} else if strings.Contains(content, `<span class="kaibai">`) {
	//	pattern := `<span class="kaibai">[a-z,A-z,0-9]+<\/span>`
	//	re := regexp.MustCompile(pattern)
	//	code := re.FindString(content)
	//	code = strings.TrimSuffix(strings.TrimPrefix(code, `<span class="kaibai">`), `</span>`)
	//	code = fmt.Sprintf("https://sexy.1433.top/%s?token=%s", code, c.token)
	//	content = re.ReplaceAllString(content, code)
	//}

	c.showMsg(msg)
}

func (c *Core) addCache(msg *WsMsgReply) {
	c.cache = append(c.cache, msg)
	if len(c.cache) >= c.cacheNum {
		removeNum := len(c.cache) - c.cacheNum
		c.cache = c.cache[removeNum:]
	}
}

func (c *Core) filterMessage(msg *WsMsgReply) {
	//if msg.IsRedPacketMsg() {
	//	switch msg.RedPackageInfo.Type {
	//	case RedPacketTypeSpecify:
	//		h.own = msg
	//	case RedPacketTypeRockPaperScissors:
	//		h.gesture = msg
	//	case RedPacketTypeHeartbeat:
	//		h.heartbeat = msg
	//	default:
	//		h.red = msg
	//	}
	//	return
	//}

	if msg.Type == WsMsgTypeMsg {
		c.addCache(msg)

		if msg.UserName == c.sdk.username {
			c.lastest = msg
		}
	}
}

func (c *Core) HandleWsStatusMsg(data interface{}) {
	msg, ok := data.(*WsMsgReply)
	if !ok {
		return
	}
	c.showMsg(msg)
}

func (c *Core) KeepLive() {
	ticker := time.NewTicker(3 * time.Minute)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case _ = <-ticker.C:
				c.eh.Pub(eventHandler.WsSend, "-hb-")
			}
		}
	}()
}

func (c *Core) showMsg(msg *WsMsgReply) {
	if c.msgChannel == nil {
		c.showMsgCache = append(c.showMsgCache, msg)
		return
	}
	if len(c.showMsgCache) != 0 {
		for _, v := range c.showMsgCache {
			c.msgChannel <- v
		}
		c.showMsgCache = nil
	}
	c.msgChannel <- msg
}

func (c *Core) ShowMsgChannel() <-chan *WsMsgReply {
	if c.msgChannel == nil {
		c.msgChannel = make(chan *WsMsgReply, 1024)
	}
	return c.msgChannel
}
