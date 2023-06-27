package core

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"fishpi/logger"
)

type Sdk struct {
	api      *Api
	ua       string
	apiKey   string
	username string

	logger logger.Logger
}

func NewSdk(api *Api, userAgent, apiKey, username string, logger logger.Logger) *Sdk {
	c := &Sdk{
		api: api,
		ua:  userAgent,

		apiKey:   apiKey,
		username: username,
		logger:   logger,
	}

	return c
}

func (c *Sdk) GetWsUrl() string {
	return c.api.wss() + "?apiKey=" + c.apiKey
}

// User 获取自己的信息
func (c *Sdk) User() (string, error) {
	body, err := c.get(c.api.user())
	if err != nil {
		return "", err
	}

	var reply userReply
	if err = json.Unmarshal(body, &reply); err != nil {
		return "", err
	}
	if reply.Code != 0 {
		return "", errors.New(reply.Msg)
	}

	return reply.Data.UserNickname, nil
}

// UserCheckedIn 获取用户是否签到 {"checkedIn":true}
func (c *Sdk) UserCheckedIn() (bool, error) {
	body, err := c.get(c.api.userCheckedIn())
	if err != nil {
		return false, err
	}
	var reply userCheckedInReply
	if err = json.Unmarshal(body, &reply); err != nil {
		return false, err
	}

	return reply.CheckedIn, nil
}

// UserLiveness 获取用户活跃度 {"liveness":87}
func (c *Sdk) UserLiveness() (float64, error) {
	body, err := c.get(c.api.userLiveness())
	if err != nil {
		return 0, err
	}

	var reply userLivenessReply
	if err = json.Unmarshal(body, &reply); err != nil {
		return 0, err
	}

	return reply.Liveness, nil
}

// ChatRecordPage 获取消息历史记录按页数
func (c *Sdk) ChatRecordPage(page int) ([]*ChatRecordPageData, error) {
	body, err := c.get(c.api.chatRecordPage(page))
	if err != nil {
		return nil, err
	}

	var reply ChatRecordPageReply
	if err = json.Unmarshal(body, &reply); err != nil {
		return nil, err
	}
	if reply.Code != 0 {
		return nil, errors.New(reply.Msg)
	}

	return reply.Data, nil
}

// UserInfo 获取用户信息
func (c *Sdk) UserInfo(username string) string {
	body, err := c.get(c.api.userInfo(username))
	if err != nil {
		return err.Error()
	}

	var uir UserInfoReply
	if err = json.Unmarshal(body, &uir); err != nil {
		return err.Error()
	}
	uir.Parse()

	return uir.String()
}

// DrawYesterdayLivenessReward 领取昨日活跃奖励 {"sum":-1}
func (c *Sdk) DrawYesterdayLivenessReward() (string, error) {
	body, err := c.get(c.api.drawYesterdayLivenessReward())
	if err != nil {
		return "", err
	}

	var reply drawYesterdayLivenessRewardReply
	if err = json.Unmarshal(body, &reply); err != nil {
		return "", err
	}

	if reply.Sum == -1 {
		return "积分已经领取", nil
	}

	return fmt.Sprintf("领取到%d积分", reply.Sum), nil
}

// IsCollectedLiveness 查询昨日奖励领取状态 {"isCollectedYesterdayLivenessReward":true}
func (c *Sdk) IsCollectedLiveness() (bool, error) {
	body, err := c.get(c.api.isCollectedLiveness())
	if err != nil {
		return false, err
	}

	var reply isCollectdLivenessReply
	if err = json.Unmarshal(body, &reply); err != nil {
		return false, err
	}

	return reply.IsCollectedYesterdayLivenessReward, nil
}

// GetArticleInfo 获取文章信息
func (c *Sdk) GetArticleInfo(articleId string) ([]byte, error) {
	body, err := c.get(c.api.getArticleInfo(articleId))
	if err != nil {
		return nil, err
	}

	return body, nil
}

// SendMsg 发送消息 {"code":0}
func (c *Sdk) SendMsg(msg string) error {
	data := &sendMsgData{
		ApiKey:  c.apiKey,
		Content: msg,
		Client:  "Golang/v0.0.2",
	}

	body, err := c.post(c.api.sendMsg(), data)
	if err != nil {
		return err
	}

	var reply sendMsgReply
	if err = json.Unmarshal(body, &reply); err != nil {
		return err
	}
	if reply.Code != 0 {
		return fmt.Errorf("send msg error, code: %d", reply.Code)
	}

	return nil
}

// SendBreezeMoon 发送消息 {"code":0}
func (c *Sdk) SendBreezeMoon(msg string) error {
	data := &sendBreezeMoonData{
		ApiKey:            c.apiKey,
		BreezeMoonContent: msg,
	}

	body, err := c.post(c.api.sendBreezeMoon(), data)
	if err != nil {
		return err
	}

	var reply sendBreezeMoonReply
	if err = json.Unmarshal(body, &reply); err != nil {
		return err
	}
	if reply.Code != 0 {
		return fmt.Errorf("send msg error, code: %d, msg: %s", reply.Code, reply.Msg)
	}

	return nil
}

func (c *Sdk) BreezeMoonList(msg string) error {
	params := strings.Split(msg, "-")

	var e error
	var page, size int
	if len(params) >= 1 {
		if size, e = strconv.Atoi(params[0]); e != nil {
			size = 20
		}
	} else {
		size = 20
	}
	if len(params) >= 2 {
		if page, e = strconv.Atoi(params[1]); e != nil {
			page = 1
		}
	} else {
		page = 1
	}

	body, err := c.get(c.api.breezeMoonList(page, size))
	if err != nil {
		return err
	}

	var reply breezeMoonReply
	if err = json.Unmarshal(body, &reply); err != nil {
		return err
	}
	if reply.Code != 0 {
		return fmt.Errorf("get breezeMoon list error, code: %d", reply.Code)
	}
	fmt.Println(reply.String())

	return nil
}

func (c *Sdk) BreezeMoonUser(msg string) error {
	if msg == "" {
		return errors.New("用户名不能为空")
	}
	params := strings.Split(msg, "-")

	var name string
	var e error
	var page, size int
	if len(params) >= 1 {
		name = params[0]
	}
	if len(params) >= 2 {
		if size, e = strconv.Atoi(params[1]); e != nil {
			size = 20
		}
	} else {
		size = 20
	}
	if len(params) >= 3 {
		if page, e = strconv.Atoi(params[2]); e != nil {
			page = 1
		}
	} else {
		page = 1
	}

	body, err := c.get(c.api.breezeMoonUser(name, page, size))
	if err != nil {
		return err
	}

	var reply breezeMoonUserReply
	if err = json.Unmarshal(body, &reply); err != nil {
		return err
	}
	if reply.Code != 0 {
		return fmt.Errorf("get breezeMoon user %s error, code: %d", name, reply.Code)
	}
	fmt.Println(reply.String())

	return nil
}

// RevokeMsg 聊天室撤回消息
func (c *Sdk) RevokeMsg(oId string) error {
	data := &revokeMsgData{
		ApiKey: c.apiKey,
		OId:    oId,
	}

	body, err := c.delete(c.api.revokeMessage(oId), data)
	if err != nil {
		return err
	}

	var reply revokeMsgReply
	if err = json.Unmarshal(body, &reply); err != nil {
		return err
	}
	if reply.Code != 0 {
		return fmt.Errorf("revoke msg error, code: %d, msg: %s", reply.Code, reply.Msg)
	}

	return nil
}

// OpenRedPacket 打开红包
func (c *Sdk) OpenRedPacket(oId string, gesture string) (string, error) {
	data := &openRedPacketData{
		ApiKey: c.apiKey,
		OId:    oId,
	}
	switch gesture {
	case "":
	case "1":
		data.Gesture = 0
	case "2":
		data.Gesture = 1
	case "3":
		data.Gesture = 2
	default:
		data.gesture()
	}

	body, err := c.post(c.api.openRedPacket(), data)
	if err != nil {
		return "", err
	}

	var reply openRedPacketReply
	if err = json.Unmarshal(body, &reply); err != nil {
		return "", err
	}

	receiveResult := "但是没有领取到欸"
	var receiveList []string
	for _, v := range reply.Who {
		if v.UserName == c.username {
			if v.UserMoney < 0 {
				receiveResult = fmt.Sprintf("血亏 损失到了%d积分", -v.UserMoney)
			} else if v.UserMoney > 0 {
				receiveResult = fmt.Sprintf("真牛 领取到了%d积分", v.UserMoney)
			} else {
				receiveResult = "你抢了个寂寞"
			}
		}
		receiveList = append(receiveList, fmt.Sprintf("- %s %s 抢到了%d积分", v.Time, v.UserName, v.UserMoney))
	}
	//receiveResult += fmt.Sprintf("他出的%s", reply.Info.GestureName()) // 接口并未返回对方出拳 但是网页有

	result := fmt.Sprintf("你打开%s发的红包(%d/%d) %s\n 领取情况：\n%s\n\n%s", reply.Info.UserName, reply.Info.Got, reply.Info.Count, receiveResult, strings.Join(receiveList, "\n"), reply.Info.Msg)

	return result, nil
}

func (c *Sdk) GetApiKey() string {
	return c.apiKey
}

func (c *Sdk) GetKey(username string, passwordMd5 string, mfaCode string) error {
	data := &getKeyData{
		NameOrEmail:  username,
		UserPassword: passwordMd5,
		MfaCode:      mfaCode,
	}

	body, err := c.post(c.api.getKey(), data)
	if err != nil {
		return err
	}

	var reply getKeyReply
	if err = json.Unmarshal(body, &reply); err != nil {
		return err
	}
	if reply.Code != 0 {
		return errors.New(reply.Msg)
	}
	c.apiKey = reply.Key

	return nil
}

func (c *Sdk) post(u *url.URL, data interface{}) ([]byte, error) {
	param, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var req *http.Request

	if req, err = http.NewRequest(http.MethodPost, u.String(), bytes.NewReader(param)); err != nil {
		return nil, err
	}

	var body []byte
	if body, err = c.do(req); err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Sdk) delete(u *url.URL, data interface{}) ([]byte, error) {
	param, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var req *http.Request

	if req, err = http.NewRequest(http.MethodDelete, u.String(), bytes.NewReader(param)); err != nil {
		return nil, err
	}

	var body []byte
	if body, err = c.do(req); err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Sdk) get(u *url.URL) ([]byte, error) {
	q := u.Query()
	q.Add("apiKey", c.apiKey)
	u.RawQuery = q.Encode()
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	var body []byte
	if body, err = c.do(req); err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Sdk) do(req *http.Request) ([]byte, error) {
	req.Header.Set("User-Agent", c.ua)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var body []byte
	if body, err = io.ReadAll(resp.Body); err != nil {
		return nil, err
	}
	//fmt.Printf("url: %s\ncode: %d\nbody: %s\n", req.URL.String(), resp.StatusCode, string(body))

	return body, nil
}
