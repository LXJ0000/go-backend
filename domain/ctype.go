package domain

type Role int

const (
	RoleAdmin       Role = iota + 1 // 管理员
	RoleUser                        // 普通用户
	RoleVisitor                     // 游客
	RoleDisableUser                 // 封号
)

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "管理员"
	case RoleUser:
		return "普通用户"
	case RoleVisitor:
		return "游客"
	case RoleDisableUser:
		return "封号"
	default:
		return "其他"
	}
}

type LoginType int

const (
	SignQQ    LoginType = iota + 1 // QQ
	SignEmail                      // Email
	SignPhone                      // Phone
)

func (s LoginType) String() string {

	switch s {
	case SignQQ:
		return "QQ"
	case SignEmail:
		return "Email"
	case SignPhone:
		return "Phone"
	default:
		return "其他"
	}
}
