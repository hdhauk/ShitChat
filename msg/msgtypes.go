package msg

type ServerResp struct {
	TimeStamp string      `json:"timestamp"`
	Sender    string      `json:"sender"`
	Resp      string      `json:"response"`
	Content   interface{} `json:"content"`
}

type ClientReq struct {
	Request string      `json:"request"`
	Content interface{} `json:"content"`
}
