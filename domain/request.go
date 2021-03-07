package domain

// Message is the main message container
type Message struct {
	Id                    int32  `json:"id"`
	Date                  int32  `json:"date"`
	FromId                int32  `json:"from_id"`
	PeerId                int32  `json:"peer_id"`
	Out                   int32  `json:"out"`
	Text                  string `json:"text"`
	ConversationMessageId int32  `json:"conversation_message_id"`
	FwdMessages           []int  `json:"fwd_messages"`
	Important             bool   `json:"important"`
	RandomId              int32  `json:"random_id"`
	Attachments           []int  `json:"attachments"`
	IsHidden              bool   `json:"is_hidden"`
	Payload               string `json:"payload"`
}

// Some client info (unused yet)
type ClientInfo struct {
	ButtonActions  []string
	Keyboard       bool
	InlineKeyboard bool
	LandId         int8
}

// Object of the message
type Object struct {
	Message    Message `json:"message"`
	ClientInfo ClientInfo
}

//easyjson:json
type Request struct {
	Type    string `json:"type"`
	Object  Object `json:"object"`
	GroupId int32  `json:"group_id"`
	EventId string `json:"event_id"`
	Secret  string `json:"secret"`
}

func (v Request) IsKeyBoardButton() bool {
	return v.Object.Message.Payload != ``
}
