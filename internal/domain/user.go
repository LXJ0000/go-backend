package domain

import (
	"context"
	"strings"
	"time"
)

const (
	XUserID       = "x-user-id"
	UserSessionID = "user-session-id"

	BizUserLogin = "login"
)

func UserLogoutKey(ssid string) string {
	return strings.Join([]string{"user", "logout", ssid}, "_")
}

type User struct {
	Model
	UserID   int64     `json:"user_id,string" gorm:"primaryKey"`
	UserName string    `json:"user_name" gorm:"unique"`
	NickName string    `json:"nick_name" gorm:"size:256"`
	Email    string    `json:"email" gorm:"default:null;unique"`
	Password string    `json:"-" gorm:"size:256"`
	AboutMe  string    `json:"about_me" gorm:"size:256"`
	Birthday time.Time `json:"birthday" gorm:"default:null"`
	Phone    string    `json:"phone" gorm:"default:null;unique"`
	Region   string    `json:"region" gorm:"default:null"`
	Avatar   string    `json:"avatar" gorm:"size:1024"`
	//Telephone string `json:"telephone" gorm:"size:20"`
	//LoginType LoginType `json:"login_type" gorm:"size:20"`
	//Role      Role      `json:"role" gorm:"default:2"` //
}

func (User) TableName() string {
	return `user`
}

type UserRepository interface {
	Create(c context.Context, user *User) error
	GetByEmail(c context.Context, email string) (User, error)
	GetByID(c context.Context, id int64) (User, error)
	GetByPhone(c context.Context, phone string) (User, error)
	FindByUserIDs(c context.Context, userIDs []int64, page, size int) ([]User, error)
	InvalidToken(c context.Context, ssid string, exp time.Duration) error
	Update(c context.Context, id int64, user *User) error
	Search(c context.Context, keyword string, page, size int) ([]User, int, error)
}

type UserUsecase interface {
	GetProfileByID(c context.Context, userID int64) (*Profile, error)
	UpdateProfile(c context.Context, userID int64, user *User) error
	Logout(c context.Context, SSID string, tokenExpiry time.Duration) error

	GetUserByEmail(c context.Context, email string) (User, error)
	CreateAccessToken(user User, ssid string, secret string, expiry int) (accessToken string, err error)
	CreateRefreshToken(user User, ssid string, secret string, expiry int) (refreshToken string, err error)

	Create(c context.Context, user *User) error

	GetUserByPhone(c context.Context, phone string) (User, error)
	Search(c context.Context, keyword string, page, size int) ([]User, int, error)
}

type Profile struct {
	UserName     string       `json:"user_name"`
	NickName     string       `json:"nick_name"`
	Email        string       `json:"email"`
	AboutMe      string       `json:"about_me"`
	Birthday     time.Time    `json:"birthday"`
	Avatar       string       `json:"avatar"`
	RelationStat RelationStat `json:"relation_stat"`
	PostCnt      int64        `json:"post_cnt"`
}

type LoginReq struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type LoginByPhoneReq struct {
	Phone    string `form:"phone" json:"phone" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
type LoginResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SignupReq struct {
	UserName string `form:"user_name" json:"user_name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

type SignupResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type SendSMSCodeReq struct {
	Phone string `form:"phone" json:"phone" binding:"required"`
}

type LoginBySmsReq struct {
	Phone string `form:"phone" json:"phone" binding:"required"`
	Code  string `form:"code" json:"code" binding:"required"`
}

type UserSearchReq struct {
	Page    int    `form:"page" json:"page"`
	Size    int    `form:"size" json:"size"`
	Keyword string `form:"keyword" json:"keyword" binding:"required"`
}

//type Role int
//
//const (
//	RoleAdmin       Role = iota + 1 // 管理员
//	RoleUser                        // 普通用户
//	RoleVisitor                     // 游客
//	RoleDisableUser                 // 封号
//)
//
//func (r Role) String() string {
//	switch r {
//	case RoleAdmin:
//		return "管理员"
//	case RoleUser:
//		return "普通用户"
//	case RoleVisitor:
//		return "游客"
//	case RoleDisableUser:
//		return "封号"
//	default:
//		return "其他"
//	}
//}
//
//type LoginType int
//
//const (
//	SignQQ    LoginType = iota + 1 // QQ
//	SignEmail                      // Email
//	SignPhone                      // Phone
//)
//
//func (s LoginType) String() string {
//
//	switch s {
//	case SignQQ:
//		return "QQ"
//	case SignEmail:
//		return "Email"
//	case SignPhone:
//		return "Phone"
//	default:
//		return "其他"
//	}
//}
