package domain

type (
	Error struct {
		ErrorCode    int    `json:"error_code"`
		ErrorMessage string `json:"error_message"`
	}
	ApiResponse struct {
		Response []VkUser `json:"response"`
		Error    Error    `json:"error"`
	}
)
