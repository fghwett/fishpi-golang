package simple

import (
	"fishpi/core"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	pagePublicChatroom = "page-public-chatroom" // 公共聊天室
	pageUserChatroom   = "page-user-chatroom"   // 私聊
	pageMoonList       = "page-moon-list"       // 明月清风
	pageIceGame        = "page-ice-game"        // 小冰游戏
)

type Simple struct {
	// UI组件
	app    *tview.Application
	layout *tview.Grid
	pages  *tview.Pages

	// 内部数据
	publicMessageChan chan []byte
	userMessageChan   chan []byte
	iceMessageChan    chan []byte

	core *core.Core
}

func NewSimple(core *core.Core) *Simple {
	pages := tview.NewPages()

	u := &Simple{
		app:   tview.NewApplication(),
		pages: pages,

		publicMessageChan: make(chan []byte, 1024),
		userMessageChan:   make(chan []byte, 1024),
		iceMessageChan:    make(chan []byte, 1024),

		core: core,
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
	list.SetSelectedBackgroundColor(tcell.NewRGBColor(150, 150, 150))

	u.layout = tview.NewGrid().
		SetColumns(10, 0)
	u.layout.SetBackgroundColor(tcell.ColorDefault)

	u.layout.AddItem(list, 0, 0, 1, 1, 0, 0, true)

	u.layout.AddItem(u.pages, 0, 1, 1, 1, 0, 0, false)
}

func (u *Simple) addPublicChatroom() {

	// 聊天框
	messageView := tview.NewTextView().SetText("暂无消息").SetChangedFunc(func() {
		u.app.Draw()
	})
	messageView.ScrollToEnd()
	messageView.Box.SetTitle(" 暂无 ")
	messageView.SetBackgroundColor(tcell.ColorDefault)
	messageView.SetTextColor(tcell.ColorDefault)
	go u.updateMessageView(messageView)

	// 信息框
	infoView := tview.NewTextView().SetText("暂无内容").SetChangedFunc(func() {
		u.app.Draw()
	})
	infoView.SetBackgroundColor(tcell.ColorDefault)
	infoView.SetTextColor(tcell.ColorDefault)

	// 输入框
	inputView := tview.NewInputField()
	inputView.SetPlaceholder(" 这里输入你要发送的消息")

	style := tcell.StyleDefault
	style = style.Background(tcell.NewRGBColor(43, 43, 43))
	inputView.SetPlaceholderStyle(style)
	inputView.SetFieldStyle(style)

	inputView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			u.core.SendPublicMsg(inputView.GetText())
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

func (u *Simple) updateMessageView(messageView *tview.TextView) {
	for {
		select {
		case msg := <-u.publicMessageChan:
			if _, err := fmt.Fprintf(messageView, "\n%s", msg); err != nil {
				panic(err)
			}
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
