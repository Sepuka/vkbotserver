package domain

type Message struct {
	Id                    int32
	Date                  int32
	FromId                int32 `json:"from_id"`
	PeerId                int32 `json:"peer_id"`
	Out                   int32
	Text                  string
	ConversationMessageId int32 `json:"conversation_message_id"`
	FwdMessages           []int
	Important             bool
	RandomId              int32 `json:"random_id"`
	Attachments           []int
	IsHidden              bool
	Payload               string
}

type ClientInfo struct {
	ButtonActions  []string
	Keyboard       bool
	InlineKeyboard bool
	LandId         int8
}

type Object struct {
	Message    Message
	ClientInfo ClientInfo
}

//easyjson:json
type Request struct {
	Type    string `json:"type"`
	Object  Object
	GroupId int32  `json:"group_id"`
	EventId string `json:"event_id"`
	Secret  string
}

func (v Request) IsKeyBoardButton() bool {
	return v.Object.Message.Payload != ``
}