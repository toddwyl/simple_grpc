package facade

import "fmt"

// User
// 返回用户信息对象
type User struct {
	Id         int64  `json:"id"`
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	ProfilePic string `json:"profile_pic"`
}

// User.String
func (u User) String() string {
	return fmt.Sprintf("Id:%d, Username:%s, Nickname:%s, Profile:%s", u.Id, u.Username, u.Nickname, u.ProfilePic)
}

// RegisterUserRequest 用户注册时的请求对象
type RegisterUserRequest struct {
	Username   string `form:"username"`
	Nickname   string `form:"nickname"`
	Password   string `form:"password"`
	ProfilePic string `form:"profile_pic"`
}

// LoginUserRequest 用户登录时的请求对象
type LoginUserRequest struct {
	Username string `form:"username" binding:"required,min=6,max=64"`
	Password string `form:"password" binding:"required,min=6,max=64"`
}

type LoginUserResponse struct {
	SessionID string `json:"session_id"`
}

// UpdateUserRequest 用户更新昵称和头像时的请求对象
// 更新昵称时，Username和Nickname不可为空
// 更新头像时，Username和ProfilePath不可为空
type UpdateUserRequest struct {
	Username   string `form:"username"`
	Nickname   string `form:"nickname"`
	ProfilePic string `form:"profile_pic"`
}

type GetUserRequest struct {
	SessionId string `form:"session_id"`
}

type EditUserRequest struct {
	SessionId  string `form:"session_id"`
	Nickname   string `form:"nickname"`
	ProfilePic string `form:"profile_pic"`
}

type RegisterUserResponse struct {
}

type GetUserResponse struct {
	Username   string `json:"username"`
	Nickname   string `json:"nickname"`
	ProfilePic string `json:"profile_pic"`
}

type EditUserResponse struct {
}
