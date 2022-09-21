package core

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"fishpi/logger"
)

type Client struct {
	sdk *Sdk

	logger logger.Logger
}

func NewClient(sdk *Sdk, logger logger.Logger) *Client {
	c := &Client{
		sdk:    sdk,
		logger: logger,
	}

	return c
}

func (c *Client) SendMode() {
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
	prefixInfo = "info-"
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
	if strings.HasPrefix(msg, prefixInfo) {
		name := strings.TrimPrefix(msg, prefixInfo)
		c.logger.Log(c.sdk.UserInfo(name))
		return
	}

	fmt.Println(c.sdk.SendMsg(msg))
	fmt.Println()
	fmt.Println()
	fmt.Println()
}

func (c *Client) handleHelp() {
	help := `help - 查看帮助信息
liveness - 查询当前活跃度（官方查询时间间隔建议为30s 本程序未作限制）
reward - 查询昨日活跃奖励是否已经领取并自动领取
info-{username} - 查询用户信息 {username}为想要查询的用户的用户名

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
