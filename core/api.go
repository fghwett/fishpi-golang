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

func (a *Api) wss() string {
	u := *a.u
	u.Path = "/chat-room-channel"
	return strings.ReplaceAll(u.String(), "https://", "wss://")
}
