package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"math/rand"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/tw"
)

type sendMsgData struct {
	ApiKey  string `json:"apiKey"`
	Content string `json:"content"`
	Client  string `json:"client"`
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

type pointTransferData struct {
	ApiKey   string `json:"apiKey"`
	Username string `json:"userName"`
	Amount   int    `json:"amount"`
	Memo     string `json:"memo"`
}

//type revokeMsgReply struct {
//	Code int    `json:"code"`
//	Msg  string `json:"msg"`
//}

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
	WsMsgTypeCustomMessage   = "customMessage"   // 进入离开聊天室 消息
	WsMsgTypeBarrage         = "barrager"        // 弹幕
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
	Time             string        `json:"time"`             // 发布时间
	UserName         string        `json:"userName"`         // 用户名
	UserNickname     string        `json:"userNickname"`     // 用户昵称
	UserAvatarURL    string        `json:"userAvatarURL"`    // 用户头像
	UserAvatarURL20  string        `json:"userAvatarURL20"`  // 用户头像 20px
	UserAvatarURL48  string        `json:"userAvatarURL48"`  // 用户头像 48px
	UserAvatarURL210 string        `json:"userAvatarURL210"` // 用户邮箱 210px
	SysMetal         string        `json:"sysMetal"`         // 徽章数据 json字符串
	Content          string        `json:"content"`          // 消息内容 HTML格式 如果是红包则是JSON格式
	Md               string        `json:"md"`               // 消息内容 Markdown格式，红包消息无此栏位
	SysMetalInfo     *SysMetalInfo // 徽章数据解析
	JsonInfo         *JsonInfo     // json内容解析

	// 红包领取消息
	Count   int    `json:"count"`   // 红包个数
	Got     int    `json:"got"`     // 已领取个数
	WhoGive string `json:"whoGive"` // 发送者用户名
	WhoGot  string `json:"whoGot"`  // 领取者用户名

	// 客户端
	Client string `json:"client"` // 消息客户端

	// 普通消息
	Message string `json:"message"` // 普通消息的消息内容

	// 弹幕消息
	BarrageColor   string `json:"barragerColor"`
	BarrageContent string `json:"barragerContent"`
}

func (w *WsMsgReply) Parse() {
	if w.Content == "" {
		return
	}

	if strings.HasPrefix(w.Content, `{\"`) {
		return
	}

	rp := new(JsonInfo)
	if err := json.Unmarshal([]byte(w.Content), rp); err != nil {
		slog.Info("解析JSON数据失败", slog.Any("err", err), slog.String("content", w.Content))
		return
	}
	w.JsonInfo = rp
}

func (w *WsMsgReply) IsRedPacketMsg() bool {
	rp := w.JsonInfo
	return rp != nil && rp.MsgType == JsonMsgTypeRedPacket
}

func (w *WsMsgReply) Msg() string {
	var result string
	switch w.Type {
	case WsMsgTypeOnline:
		result = fmt.Sprintf("当前话题：%s 在线人数：%d", w.Discussing, w.OnlineChatCnt)
	case WsMsgTypeCustomMessage:
		result = w.Message
	case WsMsgTypeBarrage:
		result = fmt.Sprintf("%s发送了弹幕消息：(%s)%s", w.UserNickname, w.BarrageColor, w.BarrageContent)
	case WsMsgTypeMsg:
		if rp := w.JsonInfo; rp != nil && rp.MsgType != "" {
			if rp.MsgType == JsonMsgTypeRedPacket {
				special := ""
				if rp.Type == RedPacketTypeSpecify {
					special = rp.Recivers
				}
				result = fmt.Sprintf("%s %s(%s): 我发了个%s%s 里面有%d积分(%d/%d)", w.Time[11:], w.UserNickname, w.UserName, rp.TypeName(), special, rp.Money, rp.Got, rp.Count)
			} else if rp.MsgType == JsonMsgTypeWeather {
				result = w.decodeJsonWeatherMsg()
			} else if rp.MsgType == JsonMsgTypeMusic {
				result = fmt.Sprintf("%s %s(%s): 让我们来听「%s」吧\n\t链接: %s\n\t封面: %s\n\t来源: %s", w.Time[11:], w.UserNickname, w.UserName, rp.Title, rp.Source, rp.CoverURL, rp.From)
			} else {
				result = fmt.Sprintf("%s %s(%s): 发送了未处理的JSON数据(%s)%s", w.Time[11:], w.UserNickname, w.UserName, rp.MsgType, w.Content)
			}
		} else if strings.Contains(w.Content, "https://www.lingmx.com/card/index2.html") {
			result = w.decodeWeatherMsg()
		} else if strings.Contains(w.Content, "https://www.lingmx.com/card/index.html") {
			result = w.decodeSingleWeatherMsg()
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

			result = fmt.Sprintf("%s %s(%s): %s(%s)", w.Time[11:], w.UserNickname, w.UserName, content, w.Client)

		}
	case WsMsgTypeDiscussChanged:
		result = fmt.Sprintf("话题变更：%s", w.NewDiscuss)
	case WsMsgTypeRedPacketStatus:
		result = fmt.Sprintf("%s领取了%s发的红包(%d/%d)", w.WhoGot, w.WhoGive, w.Got, w.Count)
	}
	return result
}

func (w *WsMsgReply) decodeSingleWeatherMsg() string {
	source := func() string {
		if w.Md != "" {
			return w.Md
		}
		return w.Content
	}()
	pattern1 := `<img src="https:\/\/img\.shields\.io\/badge\/.+">`
	re1 := regexp.MustCompile(pattern1)
	codeSource := re1.FindString(source)
	code := strings.TrimSuffix(strings.TrimPrefix(codeSource, `<img src="https://img.shields.io/badge/`), `">`)

	pattern2 := `<iframe.+iframe>`
	re2 := regexp.MustCompile(pattern2)
	linkSource := re2.FindString(source)

	var link string
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(linkSource))
	if err != nil {
		return source
	}
	dom.Find("iframe").Each(func(i int, s *goquery.Selection) {
		link, _ = s.Attr("src")
	})

	u, e := url.Parse(link)
	if e != nil {
		return source
	}
	month := u.Query().Get("m")
	day := u.Query().Get("d")
	wea := u.Query().Get("w")
	a := u.Query().Get("a")
	weather := fmt.Sprintf("%s月%s日, 天气: %s, 当前温度: %s ℃", month, day, wea, a)
	return strings.ReplaceAll(strings.ReplaceAll(source, codeSource, code), linkSource, weather)
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
		table := tablewriter.NewTable(buffer,
			tablewriter.WithConfig(tablewriter.Config{
				Header: tw.CellConfig{
					Alignment: tw.CellAlignment{
						Global: tw.AlignCenter,
					},
				},
				Row: tw.CellConfig{
					Alignment: tw.CellAlignment{
						Global: tw.AlignCenter,
					},
				},
			}),
		)
		table.Header(strings.Split(u.Query().Get("date"), ","))
		//table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})

		for _, v := range data {
			table.Append(v)
		}
		table.Render()
		msg = string(buffer.Bytes())
		msg += u.Query().Get("st")
	})
	return msg
}

func (w *WsMsgReply) decodeJsonWeatherMsg() string {
	//msg := `{"date":"4/16,4/17,4/18","st":"未来24小时多云","min":"15.49,17.49,20.49","msgType":"weather","t":"厦门","max":"25.41,26.49,26.49","weatherCode":"PARTLY_CLOUDY_DAY,CLOUDY,CLOUDY","type":"weather"}`
	msg := ``

	weather := w.JsonInfo

	weatherCodes := strings.Split(weather.WeatherCode, ",")

	weatherMap := map[string]string{
		"CLEAR_DAY":           "晴",
		"CLEAR_NIGHT":         "晴",
		"PARTLY_CLOUDY_DAY":   "多云 ",
		"PARTLY_CLOUDY_NIGHT": "多云",
		"CLOUDY":              "阴",
		"LIGHT_HAZE":          "轻度雾霾",
		"MODERATE_HAZE":       "中度雾霾",
		"HEAVY_HAZE":          "重度雾霾",
		"LIGHT_RAIN":          "小雨",
		"MODERATE_RAIN":       "中雨",
		"HEAVY_RAIN":          "大雨",
		"STORM_RAIN":          "暴雨",
		"FOG":                 "雾",
		"LIGHT_SNOW":          "小雪",
		"MODERATE_SNOW":       "中雪",
		"HEAVY_SNOW":          "大雪",
		"STORM_SNOW":          "暴雪",
		"DUST":                "浮尘",
		"SAND":                "沙尘",
		"WIND":                "大风",
	}

	var weatherWords []string
	for _, v := range weatherCodes {
		var str string
		if key, ok := weatherMap[v]; ok {
			str = key
		} else {
			str = v
		}
		weatherWords = append(weatherWords, str)
	}

	data := [][]string{
		weatherWords,
		strings.Split(weather.Max, ","),
		strings.Split(weather.Min, ","),
	}

	msg = fmt.Sprintf("%s天气\n", weather.T)
	buffer := bytes.NewBufferString(msg)
	table := tablewriter.NewTable(buffer,
		tablewriter.WithConfig(tablewriter.Config{
			Header: tw.CellConfig{
				Alignment: tw.CellAlignment{
					Global: tw.AlignCenter,
				},
			},
			Row: tw.CellConfig{
				Alignment: tw.CellAlignment{
					Global: tw.AlignCenter,
				},
			},
		}),
	)
	table.Header(strings.Split(weather.Date, ","))
	//table.SetColumnAlignment([]int{tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
	msg = string(buffer.Bytes())
	msg += weather.St

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
	JsonMsgTypeRedPacket = "redPacket"
	JsonMsgTypeWeather   = "weather"
	JsonMsgTypeMusic     = "music"

	RedPacketTypeRandom            = "random"            // 拼手气红包
	RedPacketTypeAverage           = "average"           // 平分红包
	RedPacketTypeSpecify           = "specify"           // 专属红包
	RedPacketTypeHeartbeat         = "heartbeat"         // 心跳红包
	RedPacketTypeRockPaperScissors = "rockPaperScissors" // 猜拳红包
)

// JsonInfo json的数据结构
type JsonInfo struct {
	// 消息类型 redPacket-红包 weather-天气 music-音乐
	MsgType string `json:"msgType"`
	/*
		红包类型 random(拼手气红包), average(平分红包)，specify(专属红包)，heartbeat(心跳红包)，rockPaperScissors(猜拳红包)
		天气类型 weather
		音乐类型 music
	*/
	Type string `json:"type"`

	Msg      string        `json:"msg"`      // 红包祝福语
	Recivers string        `json:"recivers"` // 红包接收者用户名，专属红包有效
	SenderId string        `json:"senderId"` // 发送者id
	Money    int           `json:"money"`    // 红包积分
	Count    int           `json:"count"`    // 红包个数
	Got      int           `json:"got"`      // 已领取个数
	Who      []interface{} `json:"who"`      // 已领取者信息

	Date        string `json:"date"`        // 日期
	St          string `json:"st"`          // 一句话
	Min         string `json:"min"`         // 最低温度
	T           string `json:"t"`           // 城市
	Max         string `json:"max"`         // 最高温度
	WeatherCode string `json:"weatherCode"` // 天气

	CoverURL string `json:"coverURL"` // 封面链接
	From     string `json:"from"`     // 来源
	Source   string `json:"source"`   // 音乐链接
	Title    string `json:"title"`    // 音乐名
}

func (jsonInfo *JsonInfo) TypeName() string {
	names := map[string]string{
		"random":            "拼手气红包",
		"average":           "平分红包",
		"specify":           "专属红包",
		"heartbeat":         "心跳红包",
		"rockPaperScissors": "猜拳红包",
		"weather":           "天气",
		"music":             "音乐",
	}
	name, ok := names[jsonInfo.Type]
	if ok {
		return name
	}
	return fmt.Sprintf("未处理类型(%s.%s)", jsonInfo.MsgType, jsonInfo.Type)
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

type ArticleInfoData struct {
	ArticleId string `json:"articleId"`
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
}

func (a *ArticleInfoData) Query() string {
	values := url.Values{}
	if a.PageSize != 0 {
		values.Add("size", strconv.Itoa(a.PageSize))
	}
	if a.Page != 0 {
		values.Add("p", strconv.Itoa(a.Page))
	}
	if values.Encode() != "" {
		return values.Encode()
	}
	return ""
}

type ArticleInfoReply struct {
	Msg  string `json:"msg"`  // 错误消息
	Code int    `json:"code"` // 为 0 则密钥有效，为 -1 则密钥无效
	Data struct {
		Article struct {
			ArticleCreateTime  string `json:"articleCreateTime"` //
			DiscussionViewable bool   `json:"discussionViewable"`
			ArticleToC         string `json:"articleToC"`
			ThankedCnt         int    `json:"thankedCnt"`
			ArticleComments    []struct {
				CommentNice              bool    `json:"commentNice"`              // 是否好评 false
				CommentCreateTimeStr     string  `json:"commentCreateTimeStr"`     // 评论创建时间 2022-12-08 09:43:14
				CommentAuthorId          string  `json:"commentAuthorId"`          // 评论作者ID 1656984017362
				CommentScore             float64 `json:"commentScore"`             // 评论分数 0
				CommentCreateTime        string  `json:"commentCreateTime"`        // 评论创建时间 Thu Dec 08 09:43:14 CST 2022
				CommentAuthorURL         string  `json:"commentAuthorURL"`         // 评论人链接 https://www.zqcnc.cn
				CommentVote              int     `json:"commentVote"`              // 评论投票 -1
				CommentRevisionCount     int     `json:"commentRevisionCount"`     // 评论修订数量 1
				TimeAgo                  string  `json:"timeAgo"`                  // 时间 7 个月前
				CommentOriginalCommentId string  `json:"commentOriginalCommentId"` // 评论原始评论ID 1670464120019
				SysMetal                 []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					Data        string `json:"data"`
					Attr        string `json:"attr"`
					Enabled     bool   `json:"enabled"`
				} `json:"sysMetal"`                                   // 评论人的徽章
				CommentGoodCnt     int    `json:"commentGoodCnt"`     // 评论点赞数 0
				CommentVisible     int    `json:"commentVisible"`     // 评论可见性 0
				CommentOnArticleId string `json:"commentOnArticleId"` // 评论所在文章ID 1670463550914
				RewardedCnt        int    `json:"rewardedCnt"`        // 评论奖励数 1
				CommentSharpURL    string `json:"commentSharpURL"`    // 评论锚点链接 /article/1670463550914#1670463793959
				CommentAnonymous   int    `json:"commentAnonymous"`   // 评论匿名性 0
				CommentReplyCnt    int    `json:"commentReplyCnt"`    // 评论回复数 0
				OId                string `json:"oId"`                // 评论ID 1670463793959
				CommentContent     string `json:"commentContent"`     // 评论内容
				CommentStatus      int    `json:"commentStatus"`      // 评论状态 0
				Commenter          struct {
					UserOnlineFlag                bool   `json:"userOnlineFlag"`
					OnlineMinute                  int    `json:"onlineMinute"`
					UserPointStatus               int    `json:"userPointStatus"`
					UserFollowerStatus            int    `json:"userFollowerStatus"`
					UserGuideStep                 int    `json:"userGuideStep"`
					UserOnlineStatus              int    `json:"userOnlineStatus"`
					UserCurrentCheckinStreakStart int    `json:"userCurrentCheckinStreakStart"`
					ChatRoomPictureStatus         int    `json:"chatRoomPictureStatus"`
					UserTags                      string `json:"userTags"`
					UserCommentStatus             int    `json:"userCommentStatus"`
					UserTimezone                  string `json:"userTimezone"`
					UserURL                       string `json:"userURL"`
					UserForwardPageStatus         int    `json:"userForwardPageStatus"`
					UserUAStatus                  int    `json:"userUAStatus"`
					UserIndexRedirectURL          string `json:"userIndexRedirectURL"`
					UserLatestArticleTime         int64  `json:"userLatestArticleTime"`
					UserTagCount                  int    `json:"userTagCount"`
					UserNickname                  string `json:"userNickname"`
					UserListViewMode              int    `json:"userListViewMode"`
					UserLongestCheckinStreak      int    `json:"userLongestCheckinStreak"`
					UserAvatarType                int    `json:"userAvatarType"`
					UserSubMailSendTime           int64  `json:"userSubMailSendTime"`
					UserUpdateTime                int64  `json:"userUpdateTime"`
					UserSubMailStatus             int    `json:"userSubMailStatus"`
					UserJoinPointRank             int    `json:"userJoinPointRank"`
					UserLatestLoginTime           int64  `json:"userLatestLoginTime"`
					UserAppRole                   int    `json:"userAppRole"`
					UserAvatarViewMode            int    `json:"userAvatarViewMode"`
					UserStatus                    int    `json:"userStatus"` // 用户状态 4为销号
					UserLongestCheckinStreakEnd   int    `json:"userLongestCheckinStreakEnd"`
					UserWatchingArticleStatus     int    `json:"userWatchingArticleStatus"`
					UserLatestCmtTime             int64  `json:"userLatestCmtTime"`
					UserProvince                  string `json:"userProvince"`
					UserCurrentCheckinStreak      int    `json:"userCurrentCheckinStreak"`
					UserNo                        int    `json:"userNo"`
					UserAvatarURL                 string `json:"userAvatarURL"`
					UserFollowingTagStatus        int    `json:"userFollowingTagStatus"`
					UserLanguage                  string `json:"userLanguage"`
					UserJoinUsedPointRank         int    `json:"userJoinUsedPointRank"`
					UserCurrentCheckinStreakEnd   int    `json:"userCurrentCheckinStreakEnd"`
					UserFollowingArticleStatus    int    `json:"userFollowingArticleStatus"`
					UserKeyboardShortcutsStatus   int    `json:"userKeyboardShortcutsStatus"`
					UserReplyWatchArticleStatus   int    `json:"userReplyWatchArticleStatus"`
					UserCommentViewMode           int    `json:"userCommentViewMode"`
					UserBreezemoonStatus          int    `json:"userBreezemoonStatus"`
					UserCheckinTime               int64  `json:"userCheckinTime"`
					UserUsedPoint                 int    `json:"userUsedPoint"`
					UserArticleStatus             int    `json:"userArticleStatus"`
					UserPoint                     int    `json:"userPoint"`
					UserCommentCount              int    `json:"userCommentCount"`
					UserIntro                     string `json:"userIntro"`
					UserMobileSkin                string `json:"userMobileSkin"`
					UserListPageSize              int    `json:"userListPageSize"`
					OId                           string `json:"oId"`
					UserName                      string `json:"userName"` // 用户名
					UserGeoStatus                 int    `json:"userGeoStatus"`
					UserLongestCheckinStreakStart int    `json:"userLongestCheckinStreakStart"`
					UserSkin                      string `json:"userSkin"`
					UserNotifyStatus              int    `json:"userNotifyStatus"`
					UserFollowingUserStatus       int    `json:"userFollowingUserStatus"`
					UserArticleCount              int    `json:"userArticleCount"`
					UserRole                      string `json:"userRole"`
				} `json:"commenter"`                                                                          // 评论者信息
				CommentAuthorName                 string `json:"commentAuthorName"`                           // 评论者昵称
				CommentThankCnt                   int    `json:"commentThankCnt"`                             // 评论感谢数 0
				CommentBadCnt                     int    `json:"commentBadCnt"`                               // 评论踩数 0
				Rewarded                          bool   `json:"rewarded"`                                    // 是否打赏过 false
				CommentAuthorThumbnailURL         string `json:"commentAuthorThumbnailURL"`                   // 评论者头像
				CommentAudioURL                   string `json:"commentAudioURL"`                             // 评论音频地址
				CommentQnAOffered                 int    `json:"commentQnAOffered"`                           // 评论是否被采纳 0
				CommentOriginalAuthorThumbnailURL string `json:"commentOriginalAuthorThumbnailURL,omitempty"` // 评论原作者头像
				PaginationCurrentPageNum          int    `json:"paginationCurrentPageNum,omitempty"`          // 当前页码
			} `json:"articleComments"`
			ArticleRewardPoint          int    `json:"articleRewardPoint"`
			ArticleRevisionCount        int    `json:"articleRevisionCount"`
			ArticleLatestCmtTime        string `json:"articleLatestCmtTime"`
			ArticleThumbnailURL         string `json:"articleThumbnailURL"`
			ArticleAuthorName           string `json:"articleAuthorName"`
			ArticleType                 int    `json:"articleType"`
			ArticleCreateTimeStr        string `json:"articleCreateTimeStr"`
			ArticleViewCount            int    `json:"articleViewCount"`
			ArticleCommentable          bool   `json:"articleCommentable"`
			ArticleAuthorThumbnailURL20 string `json:"articleAuthorThumbnailURL20"`
			ArticleOriginalContent      string `json:"articleOriginalContent"`
			ArticlePreviewContent       string `json:"articlePreviewContent"`
			ArticleContent              string `json:"articleContent"`
			ArticleAuthorIntro          string `json:"articleAuthorIntro"`
			ArticleCommentCount         int    `json:"articleCommentCount"`
			RewardedCnt                 int    `json:"rewardedCnt"`
			ArticleLatestCmterName      string `json:"articleLatestCmterName"`
			ArticleAnonymousView        int    `json:"articleAnonymousView"`
			CmtTimeAgo                  string `json:"cmtTimeAgo"`
			ArticleLatestCmtTimeStr     string `json:"articleLatestCmtTimeStr"`
			ArticleNiceComments         []struct {
				CommentCreateTimeStr     string  `json:"commentCreateTimeStr"`
				CommentAuthorId          string  `json:"commentAuthorId"`
				CommentScore             float64 `json:"commentScore"`
				CommentCreateTime        string  `json:"commentCreateTime"`
				CommentAuthorURL         string  `json:"commentAuthorURL"`
				CommentVote              int     `json:"commentVote"`
				TimeAgo                  string  `json:"timeAgo"`
				CommentOriginalCommentId string  `json:"commentOriginalCommentId"`
				SysMetal                 []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					Data        string `json:"data"`
					Attr        string `json:"attr"`
					Enabled     bool   `json:"enabled"`
				} `json:"sysMetal"`
				CommentGoodCnt     int    `json:"commentGoodCnt"`
				CommentVisible     int    `json:"commentVisible"`
				CommentOnArticleId string `json:"commentOnArticleId"`
				RewardedCnt        int    `json:"rewardedCnt"`
				CommentThankLabel  string `json:"commentThankLabel"`
				CommentSharpURL    string `json:"commentSharpURL"`
				CommentAnonymous   int    `json:"commentAnonymous"`
				CommentReplyCnt    int    `json:"commentReplyCnt"`
				OId                string `json:"oId"`
				CommentContent     string `json:"commentContent"`
				CommentStatus      int    `json:"commentStatus"`
				Commenter          struct {
					UserOnlineFlag                bool   `json:"userOnlineFlag"`
					OnlineMinute                  int    `json:"onlineMinute"`
					UserPointStatus               int    `json:"userPointStatus"`
					UserFollowerStatus            int    `json:"userFollowerStatus"`
					UserGuideStep                 int    `json:"userGuideStep"`
					UserOnlineStatus              int    `json:"userOnlineStatus"`
					UserCurrentCheckinStreakStart int    `json:"userCurrentCheckinStreakStart"`
					ChatRoomPictureStatus         int    `json:"chatRoomPictureStatus"`
					UserTags                      string `json:"userTags"`
					UserCommentStatus             int    `json:"userCommentStatus"`
					UserTimezone                  string `json:"userTimezone"`
					UserURL                       string `json:"userURL"`
					UserForwardPageStatus         int    `json:"userForwardPageStatus"`
					UserUAStatus                  int    `json:"userUAStatus"`
					UserIndexRedirectURL          string `json:"userIndexRedirectURL"`
					UserLatestArticleTime         int64  `json:"userLatestArticleTime"`
					UserTagCount                  int    `json:"userTagCount"`
					UserNickname                  string `json:"userNickname"`
					UserListViewMode              int    `json:"userListViewMode"`
					UserLongestCheckinStreak      int    `json:"userLongestCheckinStreak"`
					UserAvatarType                int    `json:"userAvatarType"`
					UserSubMailSendTime           int64  `json:"userSubMailSendTime"`
					UserUpdateTime                int64  `json:"userUpdateTime"`
					UserSubMailStatus             int    `json:"userSubMailStatus"`
					UserJoinPointRank             int    `json:"userJoinPointRank"`
					UserLatestLoginTime           int64  `json:"userLatestLoginTime"`
					UserAppRole                   int    `json:"userAppRole"`
					UserAvatarViewMode            int    `json:"userAvatarViewMode"`
					UserStatus                    int    `json:"userStatus"`
					UserLongestCheckinStreakEnd   int    `json:"userLongestCheckinStreakEnd"`
					UserWatchingArticleStatus     int    `json:"userWatchingArticleStatus"`
					UserLatestCmtTime             int64  `json:"userLatestCmtTime"`
					UserProvince                  string `json:"userProvince"`
					UserCurrentCheckinStreak      int    `json:"userCurrentCheckinStreak"`
					UserNo                        int    `json:"userNo"`
					UserAvatarURL                 string `json:"userAvatarURL"`
					UserFollowingTagStatus        int    `json:"userFollowingTagStatus"`
					UserLanguage                  string `json:"userLanguage"`
					UserJoinUsedPointRank         int    `json:"userJoinUsedPointRank"`
					UserCurrentCheckinStreakEnd   int    `json:"userCurrentCheckinStreakEnd"`
					UserFollowingArticleStatus    int    `json:"userFollowingArticleStatus"`
					UserKeyboardShortcutsStatus   int    `json:"userKeyboardShortcutsStatus"`
					UserReplyWatchArticleStatus   int    `json:"userReplyWatchArticleStatus"`
					UserCommentViewMode           int    `json:"userCommentViewMode"`
					UserBreezemoonStatus          int    `json:"userBreezemoonStatus"`
					UserCheckinTime               int64  `json:"userCheckinTime"`
					UserUsedPoint                 int    `json:"userUsedPoint"`
					UserArticleStatus             int    `json:"userArticleStatus"`
					UserPoint                     int    `json:"userPoint"`
					UserCommentCount              int    `json:"userCommentCount"`
					UserIntro                     string `json:"userIntro"`
					UserMobileSkin                string `json:"userMobileSkin"`
					UserListPageSize              int    `json:"userListPageSize"`
					OId                           string `json:"oId"`
					UserName                      string `json:"userName"`
					UserGeoStatus                 int    `json:"userGeoStatus"`
					UserLongestCheckinStreakStart int    `json:"userLongestCheckinStreakStart"`
					UserSkin                      string `json:"userSkin"`
					UserNotifyStatus              int    `json:"userNotifyStatus"`
					UserFollowingUserStatus       int    `json:"userFollowingUserStatus"`
					UserArticleCount              int    `json:"userArticleCount"`
					UserRole                      string `json:"userRole"`
				} `json:"commenter"`
				PaginationCurrentPageNum  int    `json:"paginationCurrentPageNum"`
				CommentAuthorName         string `json:"commentAuthorName"`
				CommentThankCnt           int    `json:"commentThankCnt"`
				CommentBadCnt             int    `json:"commentBadCnt"`
				Rewarded                  bool   `json:"rewarded"`
				CommentAuthorThumbnailURL string `json:"commentAuthorThumbnailURL"`
				CommentAudioURL           string `json:"commentAudioURL"`
				CommentQnAOffered         int    `json:"commentQnAOffered"`
			} `json:"articleNiceComments"`
			Rewarded                     bool    `json:"rewarded"`
			ArticleHeat                  int     `json:"articleHeat"`
			ArticlePerfect               int     `json:"articlePerfect"`
			ArticleAuthorThumbnailURL210 string  `json:"articleAuthorThumbnailURL210"`
			ArticlePermalink             string  `json:"articlePermalink"`
			ArticleCity                  string  `json:"articleCity"`
			ArticleShowInList            int     `json:"articleShowInList"`
			IsMyArticle                  bool    `json:"isMyArticle"`
			ArticleIP                    string  `json:"articleIP"`
			ArticleEditorType            int     `json:"articleEditorType"`
			ArticleVote                  int     `json:"articleVote"`
			ArticleRandomDouble          float64 `json:"articleRandomDouble"`
			ArticleAuthorId              string  `json:"articleAuthorId"`
			ArticleBadCnt                int     `json:"articleBadCnt"`
			ArticleAuthorURL             string  `json:"articleAuthorURL"`
			IsWatching                   bool    `json:"isWatching"`
			ArticleGoodCnt               int     `json:"articleGoodCnt"`
			ArticleQnAOfferPoint         int     `json:"articleQnAOfferPoint"`
			ArticleStickRemains          int     `json:"articleStickRemains"`
			TimeAgo                      string  `json:"timeAgo"`
			ArticleUpdateTimeStr         string  `json:"articleUpdateTimeStr"`
			Offered                      bool    `json:"offered"`
			ArticleWatchCnt              int     `json:"articleWatchCnt"`
			ArticleTitleEmoj             string  `json:"articleTitleEmoj"`
			ArticleTitleEmojUnicode      string  `json:"articleTitleEmojUnicode"`
			ArticleAudioURL              string  `json:"articleAudioURL"`
			ArticleAuthorThumbnailURL48  string  `json:"articleAuthorThumbnailURL48"`
			Thanked                      bool    `json:"thanked"`
			ArticleImg1URL               string  `json:"articleImg1URL"`
			ArticlePushOrder             int     `json:"articlePushOrder"`
			ArticleCollectCnt            int     `json:"articleCollectCnt"`
			ArticleTitle                 string  `json:"articleTitle"`
			IsFollowing                  bool    `json:"isFollowing"`
			ArticleTags                  string  `json:"articleTags"`
			OId                          string  `json:"oId"`
			ArticleStick                 int     `json:"articleStick"`
			ArticleTagObjs               []struct {
				TagShowSideAd     int     `json:"tagShowSideAd"`
				TagIconPath       string  `json:"tagIconPath"`
				TagStatus         int     `json:"tagStatus"`
				TagBadCnt         int     `json:"tagBadCnt"`
				TagRandomDouble   float64 `json:"tagRandomDouble"`
				TagTitle          string  `json:"tagTitle"`
				OId               string  `json:"oId"`
				TagURI            string  `json:"tagURI"`
				TagAd             string  `json:"tagAd"`
				TagGoodCnt        int     `json:"tagGoodCnt"`
				TagCSS            string  `json:"tagCSS"`
				TagCommentCount   int     `json:"tagCommentCount"`
				TagFollowerCount  int     `json:"tagFollowerCount"`
				TagSeoTitle       string  `json:"tagSeoTitle"`
				TagLinkCount      int     `json:"tagLinkCount"`
				TagSeoDesc        string  `json:"tagSeoDesc"`
				TagReferenceCount int     `json:"tagReferenceCount"`
				TagSeoKeywords    string  `json:"tagSeoKeywords"`
				TagDescription    string  `json:"tagDescription"`
			} `json:"articleTagObjs"`
			ArticleAnonymous     int    `json:"articleAnonymous"`
			ArticleThankCnt      int    `json:"articleThankCnt"`
			ArticleRewardContent string `json:"articleRewardContent"`
			RedditScore          int    `json:"redditScore"`
			ArticleUpdateTime    string `json:"articleUpdateTime"`
			ArticleStatus        int    `json:"articleStatus"`
			ArticleAuthor        struct {
				UserOnlineFlag                bool   `json:"userOnlineFlag"`
				OnlineMinute                  int    `json:"onlineMinute"`
				UserPointStatus               int    `json:"userPointStatus"`
				UserFollowerStatus            int    `json:"userFollowerStatus"`
				UserGuideStep                 int    `json:"userGuideStep"`
				UserOnlineStatus              int    `json:"userOnlineStatus"`
				UserCurrentCheckinStreakStart int    `json:"userCurrentCheckinStreakStart"`
				ChatRoomPictureStatus         int    `json:"chatRoomPictureStatus"`
				UserTags                      string `json:"userTags"`
				SysMetal                      []struct {
					Name        string `json:"name"`
					Description string `json:"description"`
					Data        string `json:"data"`
					Attr        string `json:"attr"`
					Enabled     bool   `json:"enabled"`
				} `json:"sysMetal"`
				UserCommentStatus             int    `json:"userCommentStatus"`
				UserTimezone                  string `json:"userTimezone"`
				UserURL                       string `json:"userURL"`
				UserForwardPageStatus         int    `json:"userForwardPageStatus"`
				UserUAStatus                  int    `json:"userUAStatus"`
				UserIndexRedirectURL          string `json:"userIndexRedirectURL"`
				UserLatestArticleTime         int64  `json:"userLatestArticleTime"`
				UserTagCount                  int    `json:"userTagCount"`
				UserNickname                  string `json:"userNickname"`
				UserListViewMode              int    `json:"userListViewMode"`
				UserLongestCheckinStreak      int    `json:"userLongestCheckinStreak"`
				UserAvatarType                int    `json:"userAvatarType"`
				UserSubMailSendTime           int64  `json:"userSubMailSendTime"`
				UserUpdateTime                int64  `json:"userUpdateTime"`
				UserSubMailStatus             int    `json:"userSubMailStatus"`
				UserJoinPointRank             int    `json:"userJoinPointRank"`
				UserLatestLoginTime           int64  `json:"userLatestLoginTime"`
				UserAppRole                   int    `json:"userAppRole"`
				UserAvatarViewMode            int    `json:"userAvatarViewMode"`
				UserStatus                    int    `json:"userStatus"`
				UserLongestCheckinStreakEnd   int    `json:"userLongestCheckinStreakEnd"`
				UserWatchingArticleStatus     int    `json:"userWatchingArticleStatus"`
				UserLatestCmtTime             int64  `json:"userLatestCmtTime"`
				UserProvince                  string `json:"userProvince"`
				UserCurrentCheckinStreak      int    `json:"userCurrentCheckinStreak"`
				UserNo                        int    `json:"userNo"`
				UserAvatarURL                 string `json:"userAvatarURL"`
				UserFollowingTagStatus        int    `json:"userFollowingTagStatus"`
				UserLanguage                  string `json:"userLanguage"`
				UserJoinUsedPointRank         int    `json:"userJoinUsedPointRank"`
				UserCurrentCheckinStreakEnd   int    `json:"userCurrentCheckinStreakEnd"`
				UserFollowingArticleStatus    int    `json:"userFollowingArticleStatus"`
				UserKeyboardShortcutsStatus   int    `json:"userKeyboardShortcutsStatus"`
				UserReplyWatchArticleStatus   int    `json:"userReplyWatchArticleStatus"`
				UserCommentViewMode           int    `json:"userCommentViewMode"`
				UserBreezemoonStatus          int    `json:"userBreezemoonStatus"`
				UserCheckinTime               int64  `json:"userCheckinTime"`
				UserUsedPoint                 int    `json:"userUsedPoint"`
				UserArticleStatus             int    `json:"userArticleStatus"`
				UserPoint                     int    `json:"userPoint"`
				UserCommentCount              int    `json:"userCommentCount"`
				UserIntro                     string `json:"userIntro"`
				UserMobileSkin                string `json:"userMobileSkin"`
				UserListPageSize              int    `json:"userListPageSize"`
				OId                           string `json:"oId"`
				UserName                      string `json:"userName"`
				UserGeoStatus                 int    `json:"userGeoStatus"`
				UserLongestCheckinStreakStart int    `json:"userLongestCheckinStreakStart"`
				UserSkin                      string `json:"userSkin"`
				UserNotifyStatus              int    `json:"userNotifyStatus"`
				UserFollowingUserStatus       int    `json:"userFollowingUserStatus"`
				UserArticleCount              int    `json:"userArticleCount"`
				UserRole                      string `json:"userRole"`
			} `json:"articleAuthor"`
		} `json:"article"`
		Pagination struct {
			PaginationPageCount int   `json:"paginationPageCount"` // 	评论分页数
			PaginationPageNums  []int `json:"paginationPageNums"`  // 	建议分页页码
		} `json:"pagination"` // 分页信息
	} `json:"data"`
}

type ChatroomNodeGetReply struct {
	Msg       string `json:"msg"`
	Code      int    `json:"code"`
	Data      string `json:"data"`
	ApiKey    string `json:"apiKey"`
	Avaliable []struct {
		Node   string `json:"node"`
		Name   string `json:"name"`
		Online int    `json:"online"`
	} `json:"avaliable"`
}
