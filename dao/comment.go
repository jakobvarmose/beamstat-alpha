package dao

type Comment struct {
	Id         int64
	Sender     string
	SenderName string
	IsReply    bool
	Subject    string
	Text       string
	Pending    bool
	Received   int64
	Body       string
	IsExtended bool
	ThreadHash string
}
