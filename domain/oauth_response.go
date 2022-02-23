package domain

type (
	OauthVkTokenResponse struct {
		Token            string `json:"access_token"`
		UserId           int    `json:"user_id"`
		ExpiresIn        int32  `json:"expires_in"`
		Email            string `json:"email"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
)
