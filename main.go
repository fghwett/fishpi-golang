package main

import (
	"flag"

	"fishpi/config"
	"fishpi/core"
	"fishpi/eventHandler"
	"fishpi/logger"
	"fishpi/ws"
)

// 💦
var (
	confPath = flag.String("conf", "./_tmp/config.yaml", "config path, default: ./_tmp/config.yml")
	login    = flag.Bool("login", false, "是否登录操作(false)")
	wsMode   = flag.Bool("ws", false, "是否接收消息模式(false)")
	message  = flag.Bool("msg", false, "是否发送消息模式(false)")
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
		hl := core.NewHandler(conf.Settings.MsgCacheNum, fishPiSdk, loger)

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

	// 发送消息模式
	if *message {
		client := core.NewClient(fishPiSdk, loger)
		client.SendMode()
	}

	// 默认输出帮助信息
	flag.PrintDefaults()
}
