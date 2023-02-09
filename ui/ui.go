package ui

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

type UI struct {
	// UI组件
	app    *tview.Application
	layout *tview.Grid
	pages  *tview.Pages

	// 内部数据
	publicMessageChan chan []byte
	userMessageChan   chan []byte
	iceMessageChan    chan []byte

	// 外部数据
	core *core.Core
}

func NewUI(core *core.Core) *UI {
	u := &UI{
		app:   tview.NewApplication(),
		pages: tview.NewPages(),

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

func (u *UI) Start() error {
	go u.handlePublicMsg()
	return u.app.SetRoot(u.layout, true).EnableMouse(true).Run()
}

func (u *UI) handlePublicMsg() {
	for {
		select {
		case msg := <-u.core.ShowMsgChannel():
			u.publicMessageChan <- msg
		}
	}
}

func (u *UI) makeUI() {
	chat := tview.NewButton("聊天室").SetSelectedFunc(func() {
		u.pages.SwitchToPage(pagePublicChatroom)
	})
	userChat := tview.NewButton("私聊").SetSelectedFunc(func() {
		u.pages.SwitchToPage(pageUserChatroom)
	})
	say := tview.NewButton("明月清风").SetSelectedFunc(func() {
		u.pages.SwitchToPage(pageMoonList)
	})
	iceGame := tview.NewButton("小冰游戏").SetSelectedFunc(func() {
		u.pages.SwitchToPage(pageIceGame)
	})

	u.layout = tview.NewGrid().
		SetRows(1, 0).
		SetColumns(10, 10, 10, 10, 0)

	u.layout.AddItem(chat, 0, 0, 1, 1, 0, 0, false)
	u.layout.AddItem(userChat, 0, 1, 1, 1, 0, 0, false)
	u.layout.AddItem(say, 0, 2, 1, 1, 0, 0, false)
	u.layout.AddItem(iceGame, 0, 3, 1, 1, 0, 0, false)

	u.layout.AddItem(u.pages, 1, 0, 1, 5, 0, 0, false)
}

func (u *UI) addPublicChatroom() {
	publicChatRoom := tview.NewGrid().
		SetRows(0, 1)

	messageView := tview.NewTextView().SetText("暂无消息").SetChangedFunc(func() {
		u.app.Draw()
	})
	messageView.Box.SetTitle(" 暂无 ").SetBorder(true).SetBorderAttributes(tcell.AttrNone)
	go u.updateMessageView(messageView)
	inputView := tview.NewInputField().SetPlaceholder("这里输入你要发送的消息")

	inputView.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			u.core.SendPublicMsg(inputView.GetText())
			inputView.SetText("")
		}
		if key == tcell.KeyEscape {
			inputView.SetText("")
		}
	})

	publicChatRoom.AddItem(messageView, 0, 0, 1, 1, 0, 0, false)
	publicChatRoom.AddItem(inputView, 1, 0, 1, 1, 0, 0, true)

	u.pages.AddPage(pagePublicChatroom, publicChatRoom, true, true)
}

func (u *UI) updateMessageView(messageView *tview.TextView) {
	for {
		select {
		case msg := <-u.publicMessageChan:
			if _, err := fmt.Fprintf(messageView, "%s \n", msg); err != nil {
				panic(err)
			}
		}
	}
}

func (u *UI) addUserChatroom() {
	u.pages.AddPage(pageUserChatroom, tview.NewBox().SetTitle(" 私聊界面 ").SetTitleAlign(tview.AlignRight).SetBorder(true), true, false)
}

func (u *UI) addMoonList() {
	u.pages.AddPage(pageMoonList, tview.NewBox().SetTitle(" 明月清风界面 ").SetTitleAlign(tview.AlignRight).SetBorder(true), true, false)
}

func (u *UI) addIceGame() {
	u.pages.AddPage(pageIceGame, tview.NewBox().SetTitle(" 小冰游戏界面 ").SetTitleAlign(tview.AlignRight).SetBorder(true), true, false)
}

func (u *UI) Stop() {
	u.app.Stop()
}
