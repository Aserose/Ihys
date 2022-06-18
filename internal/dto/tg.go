package dto

type TGUser struct {
	UserId int64
	ChatId int64
}

type Response struct {
	TGUser
	MsgId        int
	MsgText      string
	CallbackData string
	ExecCmd
}

// OnTappedFunc is a function of the action button.
type OnTappedFunc func(p Response)

// ExecCmd key is the callback data.
type ExecCmd map[string]OnTappedFunc
