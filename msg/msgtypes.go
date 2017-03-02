// Package msg defines the communication protocol as defined in the project text.
package msg

// ServerResp defines all messages from the server to the client.
type ServerResp struct {
	TimeStamp string      `json:"timestamp"`
	Sender    string      `json:"sender"`
	Resp      string      `json:"response"`
	Content   interface{} `json:"content"`
}

// ClientReq defines all messages from a client to the server.
type ClientReq struct {
	Request string      `json:"request"`
	Content interface{} `json:"content"`
}
