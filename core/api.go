package core

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type Api struct {
	u *url.URL
}

func NewApi(base string) (*Api, error) {
	u, e := url.Parse(base)
	if e != nil {
		return nil, e
	}

	a := &Api{u: u}
	return a, nil
}

func (a *Api) getKey() *url.URL {
	u := *a.u
	u.Path = "/api/getKey"
	return &u
}

func (a *Api) user() *url.URL {
	u := *a.u
	u.Path = "/api/user"
	return &u
}

func (a *Api) userInfo(username string) *url.URL {
	u := *a.u
	u.Path = fmt.Sprintf("/user/%s", username)
	return &u
}

func (a *Api) revokeMessage(oId string) *url.URL {
	u := *a.u
	u.Path = fmt.Sprintf("/chat-room/revoke/%s", oId)
	return &u
}

func (a *Api) pointTransfer() *url.URL {
	u := *a.u
	u.Path = "/point/transfer"
	return &u
}

func (a *Api) userCheckedIn() *url.URL {
	u := *a.u
	u.Path = "/user/checkedIn"
	return &u
}

func (a *Api) chatRecordPage(page int) *url.URL {
	u := *a.u
	u.Path = "/chat-room/more"
	value := u.Query()
	value.Add("page", strconv.Itoa(page))
	u.RawQuery = value.Encode()
	return &u
}

func (a *Api) userLiveness() *url.URL {
	u := *a.u
	u.Path = "/user/liveness"
	return &u
}

func (a *Api) openRedPacket() *url.URL {
	u := *a.u
	u.Path = "/chat-room/red-packet/open"
	return &u
}

func (a *Api) drawYesterdayLivenessReward() *url.URL {
	u := *a.u
	u.Path = "/activity/yesterday-liveness-reward-api"
	return &u
}

func (a *Api) isCollectedLiveness() *url.URL {
	u := *a.u
	u.Path = "/api/activity/is-collected-liveness"
	return &u
}

func (a *Api) getArticleInfo(data *ArticleInfoData) *url.URL {
	u := *a.u
	u.Path = fmt.Sprintf("/api/article/%s", data.ArticleId)
	u.RawQuery = data.Query()
	return &u
}

func (a *Api) sendMsg() *url.URL {
	u := *a.u
	u.Path = "/chat-room/send"
	return &u
}

func (a *Api) sendBreezeMoon() *url.URL {
	u := *a.u
	u.Path = "/breezemoon"
	return &u
}

func (a *Api) breezeMoonList(page, size int) *url.URL {
	u := *a.u
	u.Path = "/api/breezemoons"
	value := u.Query()
	value.Add("p", strconv.Itoa(page))
	value.Add("size", strconv.Itoa(size))
	u.RawQuery = value.Encode()
	return &u
}

func (a *Api) breezeMoonUser(username string, page, size int) *url.URL {
	u := *a.u
	u.Path = fmt.Sprintf("/api/user/%s/breezemoons", username)
	value := u.Query()
	value.Add("p", strconv.Itoa(page))
	value.Add("size", strconv.Itoa(size))
	u.RawQuery = value.Encode()
	return &u
}

// 聊天室连接
func (a *Api) wss() string {
	u := *a.u
	u.Path = "/chat-room-channel"
	return strings.ReplaceAll(u.String(), "https://", "wss://")
}

// 私聊连接
func (a *Api) chatChanel() string {
	u := *a.u
	u.Path = "/chat-channel"
	return strings.ReplaceAll(u.String(), "https://", "wss://")
}

// 用户通知连接 需要cookie
func (a *Api) userChannel() string {
	u := *a.u
	u.Path = "/user-channel"
	return strings.ReplaceAll(u.String(), "https://", "wss://")
}
