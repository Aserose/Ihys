package datastruct

type PlaylistItems struct {
	From  string
	Items []PlaylistItem
}

type PlaylistItem struct {
	ID      int
	OwnerId int
	Title   string
}

type AudioItems struct {
	From  string
	Items []AudioItem
}

type AudioItem struct {
	Artist string
	Title  string
	Url    string
}

type Datum struct {
	UserId int64
	ChatID int64
	MsgId  int
	Data   string
}

type ExecParam struct {
	ChatId int64
	MsgId  int
}
