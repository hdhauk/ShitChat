package msg

type ServerResp struct {
	TimeStamp string `json:"timestamp"`
	Sender    string `json:"sender"`
	Resp      string `json:"response"`
	Content   string `json:"content"`
}

type ClientReq struct {
	Request string `json:"request"`
	Content string `json:"content"`
}
