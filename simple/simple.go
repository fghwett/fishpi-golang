package simple

import (
	"fishpi/core"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	pagePublicChatroom = "page-public-chatroom" // 公共聊天室
	pageUserChatroom   = "page-user-chatroom"   // 私聊
	pageMoonList       = "page-moon-list"       // 明月清风
	pageIceGame        = "page-ice-game"        // 小冰游戏
	pageMessageMenu    = "page-message-menu"    // 信息菜单

	actionMenu            = "action_menu"             // 打开菜单
	actionRedPacket       = "action_red-packet"       // 打开红包 普通红包 平分红包 心跳红包
	actionGestureStone    = "action_gesture_stone"    // 打开猜拳红包 石头
	actionGestureScissors = "action_gesture_scissors" // 打开猜拳红包 剪刀
	actionGestureCloth    = "action_gesture_cloth"    // 打开猜拳红包 布
	actionGestureRand     = "action_gesture_rand"     // 打开猜拳红包 赌一把
)

type record struct {
	action string

	msg *core.WsMsgReply
}

type Simple struct {
	// UI组件
	app         *tview.Application
	layout      *tview.Grid
	pages       *tview.Pages
	messageView *tview.TextView
	infoView    *tview.TextView

	// 内部数据
	publicMessageChan chan *core.WsMsgReply
	userMessageChan   chan []byte
	iceMessageChan    chan []byte
	index             map[string]*record // uuid *record
	indexMu           sync.Mutex

	// 外部包
	core *core.Core
}

func NewSimple(c *core.Core) *Simple {
	pages := tview.NewPages()

	u := &Simple{
		app:   tview.NewApplication(),
		pages: pages,

		publicMessageChan: make(chan *core.WsMsgReply, 1024),
		userMessageChan:   make(chan []byte, 1024),
		iceMessageChan:    make(chan []byte, 1024),
		index:             make(map[string]*record),

		core: c,
	}

	u.addPublicChatroom()
	u.addUserChatroom()
	u.addMoonList()
	u.addIceGame()
	u.makeUI()

	return u
}

func (u *Simple) Start() error {
	go u.handlePublicMsg()
	return u.app.SetRoot(u.layout, true).EnableMouse(true).Run()
}

func (u *Simple) handlePublicMsg() {
	for {
		select {
		case msg := <-u.core.ShowMsgChannel():
			u.publicMessageChan <- msg
		}
	}
}

func (u *Simple) makeUI() {
	list := tview.NewList()
	list.AddItem("聊天室", "", 0, func() {
		u.pages.SwitchToPage(pagePublicChatroom)
	})
	list.AddItem("私聊", "", 0, func() {
		u.pages.SwitchToPage(pageUserChatroom)
	})
	list.AddItem("明月清风", "", 0, func() {
		u.pages.SwitchToPage(pageMoonList)
	})
	list.AddItem("小冰游戏", "", 0, func() {
		u.pages.SwitchToPage(pageIceGame)
	})
	list.SetCurrentItem(0)
	list.SetMainTextStyle(tcell.StyleDefault)
	list.SetBackgroundColor(tcell.ColorDefault)
	list.SetMainTextColor(tcell.NewRGBColor(191, 191, 191))
	list.SetSelectedBackgroundColor(tcell.NewRGBColor(150, 150, 150))

	u.layout = tview.NewGrid().
		SetColumns(10, 0)
	u.layout.SetBackgroundColor(tcell.ColorDefault)

	u.layout.AddItem(list, 0, 0, 1, 1, 0, 0, true)

	u.layout.AddItem(u.pages, 0, 1, 1, 1, 0, 0, false)
}

func (u *Simple) addPublicChatroom() {

	// 聊天框
	messageView := tview.NewTextView().
		SetText("暂无消息").
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			u.app.Draw()
		})
	messageView.ScrollToEnd()
	messageView.Box.SetTitle(" 暂无 ")
	messageView.SetBackgroundColor(tcell.ColorDefault)
	messageView.SetTextColor(tcell.ColorDefault)
	messageView.SetHighlightedFunc(func(added, removed, remaining []string) {
		if len(added) == 0 {
			u.showInfo(fmt.Sprintf("[Handle Message] no added! added:%s removed:%s remaining:%s", strings.Join(added, ","), strings.Join(removed, ","), strings.Join(remaining, ",")))
			return
		}
		uid := added[0]
		rec, b := u.index[uid]
		if !b {
			u.showInfo(fmt.Sprintf("[Handle Message] index uid not find: %s", uid))
			return
		}
		switch rec.action {
		case actionMenu:
			u.openMessageMenu(rec.msg)
		case actionRedPacket:
			u.openRedPacket(rec.msg, "")
		case actionGestureStone:
			u.openRedPacket(rec.msg, "1")
		case actionGestureScissors:
			u.openRedPacket(rec.msg, "2")
		case actionGestureCloth:
			u.openRedPacket(rec.msg, "3")
		case actionGestureRand:
			u.openRedPacket(rec.msg, "0")
		default:
			u.showInfo("[Handle Message] action no handle function:" + rec.action)
		}
	})
	u.messageView = messageView
	go u.updateView()

	// 信息框
	infoView := tview.NewTextView().SetText("暂无内容").SetChangedFunc(func() {
		u.app.Draw()
	})
	infoView.SetBackgroundColor(tcell.ColorDefault)
	infoView.SetTextColor(tcell.NewRGBColor(191, 191, 191))
	u.infoView = infoView

	// 输入框
	inputView := tview.NewInputField()
	inputView.SetPlaceholder(" 这里输入你要发送的消息")

	style := tcell.StyleDefault
	style = style.Background(tcell.NewRGBColor(43, 43, 43))
	style = style.Foreground(tcell.NewRGBColor(191, 191, 191))
	inputView.SetPlaceholderStyle(style)
	inputView.SetFieldStyle(style)

	inputView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			if err := u.core.SendPublicMsg(inputView.GetText()); err != nil {
				u.showInfo(fmt.Sprintf("send %s error: %s", inputView.GetText(), err))
			}
			inputView.SetText("")
		}
		if key == tcell.KeyEscape {
			inputView.SetText("")
		}
	})

	// 布局
	publicChatRoom := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(0, 80)
	publicChatRoom.SetBackgroundColor(tcell.ColorDefault)

	publicChatRoom.AddItem(messageView, 0, 0, 1, 1, 0, 0, false)
	publicChatRoom.AddItem(infoView, 0, 1, 1, 1, 0, 0, false)
	publicChatRoom.AddItem(inputView, 1, 0, 1, 2, 0, 0, false)

	u.pages.AddPage(pagePublicChatroom, publicChatRoom, true, true)
}

func (u *Simple) updateView() {
	for {
		select {
		case msg := <-u.publicMessageChan:
			u.handleMsg(msg)
		}
	}
}

func (u *Simple) addUserChatroom() {
	u.pages.AddPage(pageUserChatroom, tview.NewBox().SetTitle(" 私聊界面 ").SetTitleAlign(tview.AlignRight).SetBorder(true), true, false)
}

func (u *Simple) addMoonList() {
	u.pages.AddPage(pageMoonList, tview.NewBox().SetTitle(" 明月清风界面 ").SetTitleAlign(tview.AlignRight).SetBorder(true), true, false)
}

func (u *Simple) addIceGame() {
	u.pages.AddPage(pageIceGame, tview.NewBox().SetTitle(" 小冰游戏界面 ").SetTitleAlign(tview.AlignRight).SetBorder(true), true, false)
}

func (u *Simple) Stop() {
	u.app.Stop()
}

func (u *Simple) handleMsg(msg *core.WsMsgReply) {
	var message string
	switch msg.Type {
	case core.WsMsgTypeOnline:
		u.showInfo(fmt.Sprintf("%s: 当前话题：%s 在线人数：%d", time.Now().Format("15:04:05"), msg.Discussing, msg.OnlineChatCnt))
	case core.WsMsgTypeCustomMessage:
		message = fmt.Sprintf(`[#bfbfbf]%s`, msg.Message)
	case core.WsMsgTypeBarrage:
		// todo 处理弹幕颜色 rgba(255,255,255,1)
		message = fmt.Sprintf(`[#bbbbbb]%s发送了弹幕消息：[#bfbfbf](%s)[#bbbbbb]%s`, msg.UserNickname, msg.BarrageColor, msg.BarrageContent)
	case core.WsMsgTypeMsg:
		if rp := msg.RedPackageInfo; rp != nil && rp.Type != "" {
			special := ""
			if rp.Type == core.RedPacketTypeSpecify {
				special = rp.Recivers
			}
			action := ""
			if rp.Type == core.RedPacketTypeRockPaperScissors {
				uid1 := u.addMessageRecord(msg, actionGestureStone)
				uid2 := u.addMessageRecord(msg, actionGestureScissors)
				uid3 := u.addMessageRecord(msg, actionGestureCloth)
				rand := u.addMessageRecord(msg, actionGestureRand)
				action = fmt.Sprintf(`[#ff0000]["%s"]石头[""] ["%s"]剪刀[""] ["%s"]布[""] ["%s"]随机[""]`, uid1, uid2, uid3, rand)
			} else {
				uid := u.addMessageRecord(msg, actionRedPacket)
				action = fmt.Sprintf(fmt.Sprintf(`[#ff0000]["%s"]打开[""]`, uid))
			}
			message = fmt.Sprintf("[#bfbfbf]%s [#bbbbbb]%s[#bfbfbf](%s)[#bbbbbb]: 我发了个[#ff0000]%s%s [#bbbbbb]里面有[#ff0000]%d[#bbbbbb]积分(%d/%d) %s", msg.Time[11:], msg.UserNickname, msg.UserName, rp.TypeName(), special, rp.Money, rp.Got, rp.Count, action)
		} else if strings.Contains(msg.Content, "https://www.lingmx.com/card/index2.html") {
			//message = msg.decodeWeatherMsg()
		} else if strings.Contains(msg.Content, "https://www.lingmx.com/card/index.html") {
			//message = msg.decodeSingleWeatherMsg()
		} else {
			content := func() string {
				if msg.Md != "" {
					// todo 终端解析markdown
					return msg.Md
				}
				return msg.Content
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
						strings.Contains(str, "EXP") ||
						strings.Contains(str, "<span class='IceNet-") {
						continue
					}
					// 去除<span id = 'elves'></span>
					pattern := `<span\s+[^>]*id\s*=\s*['"]([^'"]+)['"][^>]*>(.*?)</span>`
					reg := regexp.MustCompile(pattern)
					str = reg.ReplaceAllString(str, "")
					if str == "" {
						continue
					}

					ss = append(ss, str)
				}

				return strings.Join(ss, "\n")
			}(content)

			uid := u.addMessageRecord(msg, actionMenu)
			message = fmt.Sprintf(`[#bfbfbf]%s [#bbbbbb]%s[#bfbfbf]["%s"](%s)[""][#bbbbbb]: %s[#bfbfbf](%s)`, msg.Time[11:], msg.UserNickname, uid, msg.UserName, content, msg.Client)

		}
	case core.WsMsgTypeDiscussChanged:
		u.showInfo(fmt.Sprintf("%s: 话题变更：%s", time.Now().Format("15:04:05"), msg.NewDiscuss))
	case core.WsMsgTypeRedPacketStatus:
		message = fmt.Sprintf(`[#bfbfbf]%s领取了%s发的红包(%d/%d)`, msg.WhoGot, msg.WhoGive, msg.Got, msg.Count)
	}
	if message == "" {
		return
	}

	if _, err := fmt.Fprintf(u.messageView, "\n%s", message); err != nil {
		u.showInfo(fmt.Sprintf("write %s error: %s", message, err))
	}
}

func (u *Simple) showInfo(info string) {
	if u.infoView.GetText(false) == "暂无内容" {
		u.infoView.SetText("")
	}
	if _, err := fmt.Fprintf(u.infoView, "\n%s", info); err != nil {
		panic(err)
	}
}

func (u *Simple) addMessageRecord(msg *core.WsMsgReply, action string) string {
	u.indexMu.Lock()
	defer u.indexMu.Unlock()

	str, _ := uuid.NewUUID()
	uid := str.String()
	u.index[uid] = &record{
		action: action,
		msg:    msg,
	}

	return uid
}

const (
	messageMenuRepeat = "复读机"
	messageMenuBlock  = "屏蔽此人"
	messageMenuInfo   = "查询信息"
	messageMenuClose  = "关闭"
)

func (u *Simple) openMessageMenu(msg *core.WsMsgReply) {
	if u.pages.HasPage(pageMessageMenu) {
		u.pages.HidePage(pageMessageMenu).RemovePage(pageMessageMenu)
	}
	u.pages.AddPage(
		pageMessageMenu,
		tview.NewModal().
			SetText(msg.Msg()).
			SetBackgroundColor(tcell.ColorDefault).
			AddButtons([]string{messageMenuRepeat, messageMenuBlock, messageMenuInfo, messageMenuClose}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == messageMenuRepeat {
					if err := u.core.SendPublicMsg(msg.Md); err != nil {
						u.showInfo(fmt.Sprintf("send %s error: %s", msg.Md, err))
					}
				} else if buttonLabel == messageMenuBlock {
					// todo 屏蔽功能
				} else if buttonLabel == messageMenuInfo {
					u.showInfo(u.core.GetUserInfo(msg.UserName))
				} else if buttonLabel != messageMenuClose {
					u.showInfo(fmt.Sprintf("[Message Menu] message %s action %s undefined", msg.OId, buttonLabel))
				}

				u.pages.HidePage(pageMessageMenu).RemovePage(pageMessageMenu)

			}),
		false,
		true,
	)
}

func (u *Simple) openRedPacket(msg *core.WsMsgReply, gesture string) {
	result, err := u.core.OpenRedPacket(msg.OId, gesture)
	if err != nil {
		u.showInfo(fmt.Sprintf("open %s error: %s", msg.Msg(), result))
		return
	}
	u.showInfo(result)
}
