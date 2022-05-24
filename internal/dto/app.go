package dto

type ExecCmd map[string]func(chatId int64, msgId int)

type TGUser struct {
	UserId int64
	ChatID int64
}

type Executor struct {
	TGUser
	ExecCmd
}
