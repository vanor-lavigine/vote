package response

type ApiResponse struct {
	Code      int    `json:"code"`
	Data      string `json:"data"`
	ErrorCode string `json:"errorCode"`
	ErrorMsg  string `json:"ErrorMsg"`
}

const (
	UserExists       = "UserExists"
	CreateUserFailed = "CreateUserFailed"
	LoginFailed      = "LoginFailed"
	WrongPassword    = "WrongPassword"
)
