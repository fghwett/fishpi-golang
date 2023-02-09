package main

import (
	ui2 "fishpi/ui"
	"flag"

	"fishpi/config"
	"fishpi/core"
	"fishpi/elves"
	"fishpi/eventHandler"
	"fishpi/ice"
	"fishpi/logger"
	"fishpi/ws"
)

// 💦
var (
	confPath = flag.String("conf", "./_tmp/config.yaml", "config path, default: ./_tmp/config.yml")
	login    = flag.Bool("login", false, "是否登录操作(false)")
	wsMode   = flag.Bool("ws", false, "是否接收消息模式(false)")
	message  = flag.Bool("msg", false, "是否发送消息模式(false)")
	iceMode  = flag.Bool("ice", false, "是否开启小冰游戏模式(false)")
	uiMode   = flag.Bool("ui", false, "是否使用UI模式(false)")
)

func main() {
	// 解析配置信息
	flag.Parse()

	// 初始化日志程序
	loger := logger.NewConsoleLogger()

	// 读取配置文件
	conf, err := config.NewConfig(*confPath)
	if err != nil {
		loger.Logf("读取配置文件失败 \n配置文件路径：%s\n错误信息：%s", *confPath, err)
		return
	}

	// 初始化FishPi API
	var api *core.Api
	if api, err = core.NewApi(conf.FishPi.ApiBase); err != nil {
		loger.Logf("FishPi地址信息填写失败 %s", err)
		return
	}

	fishPiSdk := core.NewSdk(api, conf.FishPi.ApiBase, conf.FishPi.ApiKey, conf.FishPi.Username, loger)

	// 登录操作
	if *login {
		if err = fishPiSdk.GetKey(conf.FishPi.Username, conf.FishPi.PasswordMd5, conf.FishPi.MfaCode); err != nil {
			loger.Logf("登陆失败 %s", err)
			return
		}

		key := fishPiSdk.GetApiKey()
		if err = conf.UpdateApiKey(key); err != nil {
			loger.Logf("更新配置文件的ApiKey错误，请手动更新\n新的ApiKey：%s\n错误信息：%s", key, err)
			return
		}
		loger.Log("更新成功！")

		return
	}

	// 接收消息模式
	if *wsMode {

		// 初始化消息处理器
		hl := core.NewHandler(conf.Settings.MsgCacheNum, conf.Elves.Token, fishPiSdk, loger)

		// 初始化事件触发器
		eh := eventHandler.NewEventHandler("websocket", loger)
		eh.Sub(eventHandler.WsMsg, hl.HandleMsg)
		eh.Sub(eventHandler.WsConnected, hl.HandleWsStatusMsg)
		eh.Sub(eventHandler.WsClosed, hl.HandleWsStatusMsg)
		eh.Sub(eventHandler.WsReconnectedFail, hl.HandleWsStatusMsg)

		// 连接ws
		u := fishPiSdk.GetWsUrl()
		wsClient := ws.NewWs(u, conf.Settings.WsInterval, eh, loger)

		if err = wsClient.Start(); err != nil {
			loger.Logf("websocket连接失败 %s", err)
			return
		}
		c := hl.KeepLive()
		go hl.Watch()
		for {
			select {
			case ping := <-c:
				wsClient.Send([]byte(ping))
			}
		}
	}

	// 小冰游戏
	if *iceMode {

		// 初始化消息处理器
		hl := ice.NewCore(conf.Ice.Ck, conf.Ice.Username, conf.Ice.Uid, loger)
		hl.SetUpdateCKFunc(conf.UpdateCK)

		// 初始化事件触发器
		eh := eventHandler.NewEventHandler("websocket", loger)
		eh.Sub(eventHandler.WsMsg, hl.HandleMsg)
		eh.Sub(eventHandler.WsConnected, hl.HandleWsStatusMsg)
		eh.Sub(eventHandler.WsClosed, hl.HandleWsStatusMsg)
		eh.Sub(eventHandler.WsReconnectedFail, hl.HandleWsStatusMsg)

		// 连接ws
		wsClient := ws.NewWs(conf.Ice.Url, conf.Settings.WsInterval, eh, loger)

		if err = wsClient.Start(); err != nil {
			loger.Logf("websocket连接失败 %s", err)
			return
		}
		c := hl.KeepLive()
		go hl.Watch()
		for {
			select {
			case msg := <-c:
				wsClient.Send(msg)
			}
		}
	}

	// 发送消息模式
	if *message {

		ec := elves.NewElves(conf.FishPi.Username, conf.Elves.Token, loger)

		eh := eventHandler.NewEventHandler("default", loger)
		eh.Sub(eventHandler.ElvesStick, ec.HandleCall)

		client := core.NewClient(fishPiSdk, eh, loger)
		client.SendMode()
	}

	// UI模式
	if *uiMode {
		// 初始化事件触发器
		eh := eventHandler.NewEventHandler("public-websocket", loger)

		// 初始化公共聊天室核心逻辑
		hl := core.NewCore(conf.Settings.MsgCacheNum, conf.Elves.Token, fishPiSdk, eh)

		eh.Sub(eventHandler.WsMsg, hl.HandleMsg)
		eh.Sub(eventHandler.WsConnected, hl.HandleWsStatusMsg)
		eh.Sub(eventHandler.WsClosed, hl.HandleWsStatusMsg)
		eh.Sub(eventHandler.WsReconnectedFail, hl.HandleWsStatusMsg)

		// 连接ws
		u := fishPiSdk.GetWsUrl()
		wsClient := ws.NewWs(u, conf.Settings.WsInterval, eh, loger)
		eh.Sub(eventHandler.WsSend, wsClient.Send)

		if err = wsClient.Start(); err != nil {
			loger.Logf("websocket连接失败 %s", err)
			return
		}

		ui := ui2.NewUI(hl)
		if err = ui.Start(); err != nil {
			panic(err)
		}
	}

	// 默认输出帮助信息
	flag.PrintDefaults()
}
