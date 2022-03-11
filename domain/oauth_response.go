package domain

type (
	// OauthVkTokenResponse carries oauth token
	OauthVkTokenResponse struct {
		Token            string `json:"access_token"`
		UserId           int    `json:"user_id"`
		ExpiresIn        int32  `json:"expires_in"`
		Email            string `json:"email"`
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description"`
	}
	UsersGetResponse struct {
		Id              int    `json:"id"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		CanAccessClosed bool   `json:"can_access_closed"`
		IsClosed        bool   `json:"is_closed"`
		Error           Error  `json:"error"`
	}
)
