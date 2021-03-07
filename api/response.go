package api

type (
	Params struct {
		Key   string
		Value string
	}
	Error struct {
		Code    int32    `json:"error_code"`
		Message string   `json:"error_msg"`
		Params  []Params `json:"request_params"`
	}

	Response struct {
		Error    Error
		Response int32
	}
)
