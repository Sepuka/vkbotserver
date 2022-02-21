package domain

type (
	OauthVkTokenResponse struct {
		Token     string `json:"access_token"`
		UserId    int32  `json:"user_id"`
		ExpiresIn int32  `json:"expires_in"`
		Email     string `json:"email"`
	}
)
