package model

type UserInfoResponse struct {
	UserID   int64  `json:"userID"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}
