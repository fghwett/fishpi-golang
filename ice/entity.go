package ice

const (
	TypeAll     = "all"     // 全局消息
	TypeSetUser = "setUser" // 登录消息
	TypeHb      = "hb"      // 保活
	TypeGameMsg = "gameMsg" // 交互消息
	TypeSetCK   = "setCK"   // 设置新的Cookie
	TypeLogin   = "login"   // 登录动作
)

type ExchangeMsg struct {
	Type  string `json:"type,omitempty"`  // 动作类型
	User  string `json:"user,omitempty"`  // 用户名
	Ck    string `json:"ck,omitempty"`    // 小冰ck
	Uid   string `json:"uid,omitempty"`   // 用户Uid
	VipLv int    `json:"vipLv,omitempty"` // VIP等级 0-普通平民 1-VIP 2-超级VIP
	Msg   string `json:"msg,omitempty"`   // 消息内容
}

func (e *ExchangeMsg) Level() string {
	switch e.VipLv {
	case 0:
		return "[白嫖怪]"
	case 1:
		return "[小冰月卡用户]"
	case 2:
		return "[小冰白金月卡用户]"
	default:
		return "[你是个什么用户？]"
	}
}
