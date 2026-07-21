package model

type User struct {
	ID       uint64
	Username string
	Nickname string
	Password string
	Email    string
}

type UserInfoResponse struct {
	UserID   int64  `json:"userId"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}
