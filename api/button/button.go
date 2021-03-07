package button

import "encoding/json"

const (
	PrimaryColor   = `primary`
	SecondaryColor = `secondary`
	NegativeColor  = `negative`
	PositiveColor  = `positive`
)

// see full docs https://vk.com/dev/bots_docs_3
type (
	// any string
	Text string
	// type belongs to the set of values: text, callback, open_link, location, vkpay, open_app
	Type string
	// there're action of button & args
	Action struct {
		Type    Type   `json:"type"`
		Label   Text   `json:"label"`
		Payload string `json:"payload"`
	}

	// Color belongs to the set: primary, secondary, negative, positive
	Button struct {
		Action Action `json:"action"`
		Color  string `json:"color"`
	}

	Keyboard struct {
		OneTime bool       `json:"one_time"`
		Buttons [][]Button `json:"buttons"`
	}

	// Contains button payload like `{"command": "start"}`
	Payload struct {
		Command string `json:"command"`
		Button  string `json:"button"`
	}
)

func (o Payload) String() string {
	var str []byte
	str, _ = json.Marshal(o)

	return string(str)
}
