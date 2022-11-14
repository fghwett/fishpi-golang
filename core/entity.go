package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"
)

type sendMsgData struct {
	ApiKey  string `json:"apiKey"`
	Content string `json:"content"`
}

type sendBreezeMoonData struct {
	ApiKey            string `json:"apiKey"`
	BreezeMoonContent string `json:"breezemoonContent"`
}

type revokeMsgData struct {
	ApiKey string `json:"apiKey"`
	OId    string `json:"oId"`
}

type revokeMsgReply struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type openRedPacketData struct {
	ApiKey  string `json:"apiKey"`
	OId     string `json:"oId"`
	Gesture int    `json:"gesture"` // 0 = 石头，1 = 剪刀，2 = 布
}

func (o *openRedPacketData) gesture() {
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(10)
	if r <= 2 {
		o.Gesture = 0
	} else if r <= 6 {
		o.Gesture = 1
	} else {
		o.Gesture = 2
	}
}

type userCheckedInReply struct {
	CheckedIn bool `json:"checkedIn"`
}

type userLivenessReply struct {
	Liveness float64 `json:"liveness"`
}

type drawYesterdayLivenessRewardReply struct {
	Sum int `json:"sum"`
}

type isCollectdLivenessReply struct {
	IsCollectedYesterdayLivenessReward bool `json:"isCollectedYesterdayLivenessReward"`
}

type sendMsgReply struct {
	Code int `json:"code"`
}

type sendBreezeMoonReply struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type breezeMoonReply struct {
	Code        int               `json:"code"`
	BreezeMoons []*breezeMoonInfo `json:"breezemoons"`
}

type breezeMoonInfo struct {
	BreezemoonAuthorName           string `json:"breezemoonAuthorName"`           // 发布者名称
	BreezemoonUpdated              int64  `json:"breezemoonUpdated"`              // 更新时间 13位毫秒
	OId                            string `json:"oId"`                            // 发布人Id
	BreezemoonCreated              int64  `json:"breezemoonCreated"`              // 创建时间
	BreezemoonAuthorThumbnailURL48 string `json:"breezemoonAuthorThumbnailURL48"` // 发布人头像
	TimeAgo                        string `json:"timeAgo"`                        // 时间格式化
	BreezemoonContent              string `json:"breezemoonContent"`              // 内容
	BreezemoonCreateTime           string `json:"breezemoonCreateTime"`           // 创建时间格式化
	BreezemoonCity                 string `json:"breezemoonCity"`                 // 发布地区
}

func (br *breezeMoonReply) String() string {
	sort.Slice(br.BreezeMoons, func(i, j int) bool {
		return br.BreezeMoons[i].BreezemoonCreated < br.BreezeMoons[j].BreezemoonCreated
	})
	var bi []string
	for _, b := range br.BreezeMoons {
		bi = append(bi, b.String())
	}
	return strings.Join(bi, "\n")
}

func (bi *breezeMoonInfo) String() string {
	ct := time.UnixMilli(bi.BreezemoonCreated).Format("2006-01-02 15:04:05")
	content := strings.TrimPrefix(strings.TrimSuffix(bi.BreezemoonContent, "</p>"), "<p>")
	return fmt.Sprintf("%s %s(%s): %s(%s)", ct, bi.BreezemoonAuthorName, bi.BreezemoonCity, content, bi.TimeAgo)
}

type breezeMoonUserReply struct {
	Code int                 `json:"code"`
	Data *breezeMoonUserData `json:"data"`
}

type breezeMoonUserData struct {
	Pagination struct {
		PaginationPageCount   int   `json:"paginationPageCount"`   // 总页数
		PaginationPageNums    []int `json:"paginationPageNums"`    // 总条数
		PaginationRecordCount int   `json:"paginationRecordCount"` // 页码
	} `json:"pagination"`
	BreezeMoons []*breezeMoonInfo `json:"breezemoons"`
}

func (br *breezeMoonUserReply) String() string {
	sort.Slice(br.Data.BreezeMoons, func(i, j int) bool {
		return br.Data.BreezeMoons[i].BreezemoonCreated < br.Data.BreezeMoons[j].BreezemoonCreated
	})
	var bi []string
	for _, b := range br.Data.BreezeMoons {
		bi = append(bi, b.String())
	}
	return strings.Join(bi, "\n")
}

type openRedPacketReply struct {
	Recivers []interface{}       `json:"recivers"`
	Who      []*openRedPacketWho `json:"who"`
	Info     *openRedPacketInfo  `json:"info"`
}

type openRedPacketWho struct {
	UserMoney int    `json:"userMoney"`
	Time      string `json:"time"`
	Avatar    string `json:"avatar"`
	UserName  string `json:"userName"`
	UserId    string `json:"userId"` // 用户id不是显示的id 而是注册时间毫秒时间戳构成的UserId
}

type openRedPacketInfo struct {
	Msg              string `json:"msg"`
	UserAvatarURL    string `json:"userAvatarURL"`
	UserAvatarURL20  string `json:"userAvatarURL20"`
	Count            int    `json:"count"`
	UserName         string `json:"userName"`
	UserAvatarURL210 string `json:"userAvatarURL210"`
	Got              int    `json:"got"`
	Gesture          int    `json:"gesture"`
	UserAvatarURL48  string `json:"userAvatarURL48"`
}

func (o *openRedPacketInfo) GestureName() string {
	switch o.Gesture {
	case 0:
		return "石头"
	case 1:
		return "剪刀"
	case 2:
		return "布"
	default:
		return strconv.Itoa(o.Gesture)
	}
}

type getKeyData struct {
	NameOrEmail  string `json:"nameOrEmail"`
	UserPassword string `json:"userPassword"`
	MfaCode      string `json:"mfaCode"`
}

type getKeyReply struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Key  string `json:"Key"`
}

type userReply struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data struct {
		UserCity           string `json:"userCity"`
		UserOnlineFlag     bool   `json:"userOnlineFlag"`
		UserPoint          int    `json:"userPoint"`
		UserAppRole        string `json:"userAppRole"`
		UserIntro          string `json:"userIntro"`
		UserNo             string `json:"userNo"`
		OnlineMinute       int    `json:"onlineMinute"`
		UserAvatarURL      string `json:"userAvatarURL"`
		UserNickname       string `json:"userNickname"`
		OId                string `json:"oId"`
		UserName           string `json:"userName"`
		CardBg             string `json:"cardBg"`
		FollowingUserCount int    `json:"followingUserCount"`
		SysMetal           string `json:"sysMetal"`
		UserRole           string `json:"userRole"`
		FollowerCount      int    `json:"followerCount"`
		UserURL            string `json:"userURL"`
	} `json:"data"`
}

const (
	WsMsgTypeOnline          = "online"          // 在线
	WsMsgTypeDiscussChanged  = "discussChanged"  // 话题变更
	WsMsgTypeRevoke          = "revoke"          // 撤回
	WsMsgTypeMsg             = "msg"             // 聊天
	WsMsgTypeRedPacketStatus = "redPacketStatus" // 红包领取
)

// WsMsgReply websocket收到的消息结构体
type WsMsgReply struct {
	Type string `json:"type"` // 消息类型

	// 在线消息
	Discussing    string            `json:"discussing"`    // 讨论的话题
	OnlineChatCnt int               `json:"onlineChatCnt"` // 在线人数
	Users         []*OnlineUserInfo `json:"users"`         // 在线用户信息

	// 话题变更
	NewDiscuss string `json:"newDiscuss"` // 新的话题内容

	// 聊天 撤回 红包领取
	OId string `json:"oId"` // 消息ID

	// 聊天消息
	Time             string         `json:"time"`             // 发布时间
	UserName         string         `json:"userName"`         // 用户名
	UserNickname     string         `json:"userNickname"`     // 用户昵称
	UserAvatarURL    string         `json:"userAvatarURL"`    // 用户头像
	UserAvatarURL20  string         `json:"userAvatarURL20"`  // 用户头像 20px
	UserAvatarURL48  string         `json:"userAvatarURL48"`  // 用户头像 48px
	UserAvatarURL210 string         `json:"userAvatarURL210"` // 用户邮箱 210px
	SysMetal         string         `json:"sysMetal"`         // 徽章数据 json字符串
	Content          string         `json:"content"`          // 消息内容 HTML格式 如果是红包则是JSON格式
	Md               string         `json:"md"`               // 消息内容 Markdown格式，红包消息无此栏位
	SysMetalInfo     *SysMetalInfo  // 徽章数据解析
	RedPackageInfo   *RedPacketInfo // 红包消息数据

	// 红包领取消息
	Count   int    `json:"count"`   // 红包个数
	Got     int    `json:"got"`     // 已领取个数
	WhoGive string `json:"whoGive"` // 发送者用户名
	WhoGot  string `json:"whoGot"`  // 领取者用户名
}

func (w *WsMsgReply) Parse() {
	if w.Content == "" {
		return
	}

	if w.Md != "" {
		return
	}

	var rp RedPacketInfo
	if err := json.Unmarshal([]byte(w.Content), &rp); err != nil {
		log.Printf("解析红包数据失败：%s, content: %s\n", err, w.Content)
		return
	}
	w.RedPackageInfo = &rp
}

func (w *WsMsgReply) IsRedPacketMsg() bool {
	rp := w.RedPackageInfo
	return rp != nil && rp.Type != ""
}

func (w *WsMsgReply) Msg() string {
	var result string
	switch w.Type {
	case WsMsgTypeOnline:
		result = fmt.Sprintf("当前话题：%s 在线人数：%d", w.Discussing, w.OnlineChatCnt)
	case WsMsgTypeMsg:
		if rp := w.RedPackageInfo; rp != nil && rp.Type != "" {
			special := ""
			if rp.Type == RedPacketTypeSpecify {
				special = rp.Recivers
			}
			result = fmt.Sprintf("%s %s(%s): 我发了个%s%s 里面有%d积分(%d/%d)", w.Time[11:], w.UserNickname, w.UserName, rp.TypeName(), special, rp.Money, rp.Got, rp.Count)
		} else if strings.Contains(w.Content, "https://www.lingmx.com/card/index2.html") {
			result = w.decodeWeatherMsg()
		} else {
			content := func() string {
				if w.Md != "" {
					return w.Md
				}
				return w.Content
			}()

			content = func(msg string) string {
				strs := strings.Split(msg, "\n")

				var ss []string
				for _, str := range strs {
					if strings.HasPrefix(str, ">") ||
						strings.HasPrefix(str, "##### 引用") ||
						strings.TrimSpace(str) == "" ||
						strings.Contains(str, "https://zsh4869.github.io/fishpi.io/?hyd=") ||
						strings.Contains(str, "extension-message") ||
						strings.Contains(str, ":sweat_drops:") ||
						strings.Contains(str, "下次更新时间") ||
						strings.Contains(str, "https://unv-shield.librian.net/api/unv_shield") ||
						strings.Contains(str, "EXP") {
						continue
					}
					ss = append(ss, str)
				}

				return strings.Join(ss, "\n")
			}(content)

			result = fmt.Sprintf("%s %s(%s): %s", w.Time[11:], w.UserNickname, w.UserName, content)

		}
	case WsMsgTypeDiscussChanged:
		result = fmt.Sprintf("话题变更：%s", w.NewDiscuss)
	case WsMsgTypeRedPacketStatus:
		result = fmt.Sprintf("%s领取了%s发的红包(%d/%d)", w.WhoGot, w.WhoGive, w.Got, w.Count)
	}
	return result
}

func (w *WsMsgReply) decodeWeatherMsg() string {
	//str := `<iframe src="https://www.lingmx.com/card/index2.html?date=8/19,8/20,8/21,8/22,8/23&weatherCode=LIGHT_RAIN,LIGHT_RAIN,CLOUDY,CLOUDY,LIGHT_RAIN&max=32,33,35,36,36&min=26,26,26,27,27&t=厦门&st=31分钟后开始下小雨，但56分钟后会停" width="380" height="370" frameborder="0"></iframe>`
	msg := w.Content
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(msg))
	if err != nil {
		return fmt.Sprintf("parse %s error: %s", w.Content, err)
	}
	dom.Find(`iframe`).Each(func(i int, s *goquery.Selection) {
		src, exist := s.Attr("src")
		if !exist {
			return
		}
		u, e := url.Parse(src)
		if e != nil {
			msg = fmt.Sprintf("parse %s error: %s", src, err)
			return
		}
		msg = u.Query().Get("t") + "天气" + "\n"
		data := [][]string{
			strings.Split(u.Query().Get("weatherCode"), ","),
			strings.Split(u.Query().Get("max"), ","),
			strings.Split(u.Query().Get("min"), ","),
		}

		buffer := bytes.NewBufferString(msg)
		table := tablewriter.NewWriter(buffer)
		table.SetHeader(strings.Split(u.Query().Get("date"), ","))
		table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})

		for _, v := range data {
			table.Append(v)
		}
		table.Render()
		msg = string(buffer.Bytes())
		msg += u.Query().Get("st")
	})
	return msg
}

// websocket 在线用户信息
type OnlineUserInfo struct {
	UserName         string `json:"userName"`         // 用户名
	HomePage         string `json:"homePage"`         // 用户首页
	UserAvatarURL    string `json:"userAvatarURL"`    // 用户头像
	UserAvatarURL20  string `json:"userAvatarURL20"`  // 用户头像 20px
	UserAvatarURL48  string `json:"userAvatarURL48"`  // 用户头像 48px
	UserAvatarURL210 string `json:"userAvatarURL210"` // 用户邮箱 210px
}

// websocket 聊天消息sysMetal
type SysMetalInfo struct {
	List []*MetalInfo `json:"list"` // 徽章列表数据
}

// websocket 徽章信息
type MetalInfo struct {
	Data        string `json:"data"`        // 徽章数据
	Name        string `json:"name"`        // 徽章名称
	Description string `json:"description"` // 徽章描述
	Attr        string `json:"attr"`        // 徽章数据，包含徽章图地址url, 背景色 backcolor, 前景色 fontcolor
	Enabled     bool   `json:"enabled"`
}

const (
	RedPacketTypeRandom            = "random"            // 拼手气红包
	RedPacketTypeAverage           = "average"           // 平分红包
	RedPacketTypeSpecify           = "specify"           // 专属红包
	RedPacketTypeHeartbeat         = "heartbeat"         // 心跳红包
	RedPacketTypeRockPaperScissors = "rockPaperScissors" // 猜拳红包
)

// websocket 红包信息解码
type RedPacketInfo struct {
	Msg      string        `json:"msg"`      // 红包祝福语
	Recivers string        `json:"recivers"` // 红包接收者用户名，专属红包有效
	SenderId string        `json:"senderId"` // 发送者id
	MsgType  string        `json:"msgType"`  // 固定 redPacket
	Money    int           `json:"money"`    // 红包积分
	Count    int           `json:"count"`    // 红包个数
	Type     string        `json:"type"`     // 红包类型 random(拼手气红包), average(平分红包)，specify(专属红包)，heartbeat(心跳红包)，rockPaperScissors(猜拳红包)
	Got      int           `json:"got"`      // 已领取个数
	Who      []interface{} `json:"who"`      // 已领取者信息
}

func (r *RedPacketInfo) TypeName() string {
	var result string
	switch r.Type {
	case RedPacketTypeRandom:
		result = "拼手气红包"
	case RedPacketTypeAverage:
		result = "平分红包"
	case RedPacketTypeSpecify:
		result = "专属红包"
	case RedPacketTypeHeartbeat:
		result = "心跳红包"
	case RedPacketTypeRockPaperScissors:
		result = "猜拳红包"
	default:
		result = "这我也不知道是什么红包" + r.Type
	}
	return result
}

type UserInfoReply struct {
	UserCity           string `json:"userCity"`           // 所在城市		1
	UserOnlineFlag     bool   `json:"userOnlineFlag"`     // 是否在线		1
	UserPoint          int    `json:"userPoint"`          // 用户积分		1
	UserAppRole        string `json:"userAppRole"`        // 用户角色 0-黑客 1-画家
	UserIntro          string `json:"userIntro"`          // 用户简介		1
	UserNo             string `json:"userNo"`             // 用户编号		1
	OnlineMinute       int    `json:"onlineMinute"`       // 在线分钟数		1
	UserAvatarURL      string `json:"userAvatarURL"`      // 用户头像链接		0
	UserNickname       string `json:"userNickname"`       // 用户昵称		1
	OId                string `json:"oId"`                // OId			0
	UserName           string `json:"userName"`           // 用户名			1
	CardBg             string `json:"cardBg"`             // 卡片背景图片		0
	AllMetalOwned      string `json:"allMetalOwned"`      // 拥有徽章		0
	FollowingUserCount int    `json:"followingUserCount"` // 关注用户数		1
	UserAvatarURL20    string `json:"userAvatarURL20"`    // 用户头像20px	0
	SysMetal           string `json:"sysMetal"`           // 徽章列表		0
	CanFollow          string `json:"canFollow"`          // 是否可以关注 no/yes/hide		0
	UserRole           string `json:"userRole"`           // 用户组名称 OP	1
	UserAvatarURL210   string `json:"userAvatarURL210"`   // 用户头像210px	0
	FollowerCount      int    `json:"followerCount"`      // 被关注数		1
	UserURL            string `json:"userURL"`            // 用户链接		1
	UserAvatarURL48    string `json:"userAvatarURL48"`    // 用户头像48px	0

	UserMetal *UserMetal // 徽章列表数据格式
}

func (u *UserInfoReply) Parse() {
	var um UserMetal
	if err := json.Unmarshal([]byte(u.SysMetal), &um); err != nil {
		log.Printf("%s的勋章信息解析失败 %s", u.UserName, err)
		return
	}
	u.UserMetal = &um
}

func (u *UserInfoReply) String() string {
	t, _ := strconv.ParseInt(u.OId, 10, 64)
	info := fmt.Sprintf("%s - %s(%s) %s\n介绍信息：%s\n链接：%s\n角色：%s(%s)\t城市：%s\n积分：%d\t在线%s\n关注数：%d\t被关注数：%d\n注册时间：%s\n",
		u.UserNo, u.UserNickname, u.UserName, u.OnlineState(),
		u.UserIntro,
		u.UserURL,
		u.UserRole, u.AppRole(), u.UserCity,
		u.UserPoint, u.OnlineTime(),
		u.FollowingUserCount, u.FollowerCount,
		time.UnixMilli(t).Format("2006-01-02 15:04:05"))
	if u.UserMetal != nil {
		var metals []string
		for _, v := range u.UserMetal.List {
			metals = append(metals, fmt.Sprintf("\t%s-%s", v.Name, v.Description))
		}
		info += "徽章列表：\n" + strings.Join(metals, "\n") + "\n"
	}
	return info
}

func (u *UserInfoReply) AppRole() string {
	if u.UserAppRole == "0" {
		return "黑客"
	}
	switch u.UserAppRole {
	case "0":
		return "黑客"
	case "1":
		return "画家"
	default:
		return u.UserAppRole
	}
}

func (u *UserInfoReply) OnlineState() string {
	if u.UserOnlineFlag {
		return "√"
	}
	return "×"
}

func (u *UserInfoReply) OnlineTime() string {
	var min, hour, day int
	min = u.OnlineMinute % 60
	hour = u.OnlineMinute / 60
	if hour < 24 {
		return fmt.Sprintf("%d小时%d分钟", hour, min)
	}

	day = hour / 24
	hour = hour % 24
	return fmt.Sprintf("%d天%d小时%d分钟", day, hour, min)
}

// UserMetal 用户勋章信息
type UserMetal struct {
	List []struct {
		Data        string `json:"data"`
		Name        string `json:"name"`        // 姓名
		Description string `json:"description"` // 描述
		Attr        string `json:"attr"`        // 背景图片
		Enabled     bool   `json:"enabled"`     // 是否启动
	} `json:"list"`
}

// ChatRecordPageReply 历史记录
type ChatRecordPageReply struct {
	Msg  string                `json:"msg"`
	Code int                   `json:"code"`
	Data []*ChatRecordPageData `json:"data"`
}

type ChatRecordPageData struct {
	UserAvatarURL    string `json:"userAvatarURL"`
	UserAvatarURL20  string `json:"userAvatarURL20"`
	UserNickname     string `json:"userNickname"`
	SysMetal         string `json:"sysMetal"`
	Time             string `json:"time"`
	OId              string `json:"oId"`
	UserName         string `json:"userName"`
	UserAvatarURL210 string `json:"userAvatarURL210"`
	Content          string `json:"content"`
	UserAvatarURL48  string `json:"userAvatarURL48"`
}
