package core

import (
	"bufio"
	"fishpi/eventHandler"
	"fmt"
	"log"
	"os"
	"strings"

	"fishpi/logger"
)

type Client struct {
	sdk *Sdk
	ln  *lnClient

	eh     eventHandler.EventHandler
	logger logger.Logger
}

func NewClient(sdk *Sdk, eh eventHandler.EventHandler, logger logger.Logger) *Client {
	c := &Client{
		sdk:    sdk,
		eh:     eh,
		logger: logger,
	}

	return c
}

func (c *Client) SendMode() {
	liveness, e := c.sdk.UserLiveness()
	if e != nil {
		liveness = 0
		fmt.Println("获取当前活跃度失败", e)
	}
	fmt.Println("当前活跃度：", liveness)
	f := func(l float64) {
		l1, e1 := c.sdk.UserLiveness()
		if e1 != nil {
			fmt.Println("获取当前活跃度失败", e1)
			return
		}
		l = l1
	}
	c.ln = NewLnClient(liveness, f)
	for {
		var buf [1024]byte
		read := bufio.NewReader(os.Stdin)
		m, err := read.Read(buf[:])
		if err != nil {
			log.Printf("read error: %s\n", err)
			continue
		}
		recv := strings.Split(string(buf[:m]), "\n")[0]
		go c.handleSendMsg(recv)
	}
}

const (
	prefixInfo           = "info-"
	prefixBreezeMoon     = "bb-"
	prefixChangeTopic    = "topic-"
	prefixBreezeMoonList = "bb-list-"
	prefixBreezeMoonUser = "bb-user-"
	prefixBarrage        = "barrage-"
)

func (c *Client) handleSendMsg(msg string) {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return
	}
	if msg == "help" {
		c.handleHelp()
		return
	}
	if msg == "liveness" {
		c.handleLiveness()
		return
	}
	if msg == "reward" {
		c.handleReward()
		return
	}
	if msg == "stick" {
		c.eh.Pub(eventHandler.ElvesStick, nil)
		return
	}
	if strings.HasPrefix(msg, prefixInfo) {
		name := strings.TrimPrefix(msg, prefixInfo)
		c.logger.Log(c.sdk.UserInfo(name))
		return
	}

	if strings.HasPrefix(msg, prefixBreezeMoonList) {
		msg = strings.TrimPrefix(msg, prefixBreezeMoonList)
		if err := c.sdk.BreezeMoonList(msg); err != nil {
			fmt.Println(err)
			return
		}
	} else if strings.HasPrefix(msg, prefixBreezeMoonUser) {
		msg = strings.TrimPrefix(msg, prefixBreezeMoonUser)
		if err := c.sdk.BreezeMoonUser(msg); err != nil {
			fmt.Println(err)
			return
		}
	} else if strings.HasPrefix(msg, prefixBreezeMoon) {
		msg = strings.TrimPrefix(msg, prefixBreezeMoon)
		if err := c.sdk.SendBreezeMoon(msg); err != nil {
			fmt.Println(err)
			return
		}
	} else {
		if strings.HasPrefix(msg, prefixChangeTopic) {
			msg = fmt.Sprintf("[setdiscuss]%s[/setdiscuss]", strings.TrimPrefix(msg, prefixChangeTopic))
		}
		if strings.HasPrefix(msg, prefixBarrage) {
			msg = strings.TrimPrefix(msg, prefixBarrage)
			color := "#66CCFF"
			msg = fmt.Sprintf(`[barrager]{"color":"%s","content":"%s"}[/barrager]`, color, msg)
		}
		if err := c.sdk.SendMsg(msg); err != nil {
			fmt.Println(err)
			return
		}
		c.ln.Say()
	}
	fmt.Println()
	fmt.Println()
	fmt.Println()
}

func (c *Client) handleHelp() {
	help := `help - 查看帮助信息
liveness - 查询当前活跃度（官方查询时间间隔建议为30s 本程序未作限制）
reward - 查询昨日活跃奖励是否已经领取并自动领取
stick - 召唤小飞棍
info-{username} - 查询用户信息 {username}为想要查询的用户的用户名
bb-{messsage} - 发布明月清风
topic-{new topic content} - 发布新话题
bb-list-{20-1} - 获取明月清风 每页20条 第一页
bb-user-{username-20-1} 获取username的明月清风 每页20条 第一页

其余信息将作为普通信息直接发送`

	c.logger.Log(help)
}

func (c *Client) handleLiveness() {
	ln, err := c.sdk.UserLiveness()
	if err != nil {
		c.logger.Logf("获取活跃度失败：%s", err)
		return
	}
	c.logger.Logf("当前活跃度：%.2f", ln)
}

func (c *Client) handleReward() {
	b, e := c.sdk.IsCollectedLiveness()
	if e != nil {
		c.logger.Logf("查询是否领取昨日活跃奖励失败 %s", e)
		return
	}
	if b {
		c.logger.Log("已经领取了昨日活跃奖励")
		return
	}
	point, err := c.sdk.DrawYesterdayLivenessReward()
	if err != nil {
		c.logger.Logf("领取昨日活跃奖励失败 %s", err)
		return
	}
	c.logger.Logf("领到%s积分", point)
}
