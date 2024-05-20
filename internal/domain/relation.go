package domain

import "context"

const (
	Follow   = true
	UnFollow = false
)

type Relation struct {
	Model
	RelationID int64 `gorm:"primaryKey;autoIncrement"`
	Followee   int64 `gorm:"not null;uniqueIndex:idx_followee_follower"` // Follower 关注了 Followee
	Follower   int64 `gorm:"not null;uniqueIndex:idx_followee_follower;index:idx_follower"`
	// 典型场景：某个人关注列表follower 某个人的粉丝列表followee 我都要
	Status bool // true 关注 false 取消关注
}

type RelationStat struct {
	Model
	UserID   int64 `gorm:"unique"`
	Follower int   // fans
	Followee int
}

type RelationUsecase interface {
	Follow(c context.Context, follower, followee int64) error
	CancelFollow(c context.Context, follower, followee int64) error
	GetFollower(c context.Context, userID int64, page, size int) ([]User, error) // 粉丝列表
	GetFollowee(c context.Context, userID int64, page, size int) ([]User, error) // 工作者列表
	Detail(c context.Context, follower, followee int64) (Relation, error)        // 关注状态
	Stat(c context.Context, userID int64) (RelationStat, error)
}

type RelationRepository interface {
	Follow(c context.Context, follower, followee int64) error
	CancelFollow(c context.Context, follower, followee int64) error
	GetFollower(c context.Context, follower int64, page, size int) ([]Relation, error) // 粉丝列表
	GetFollowee(c context.Context, follower int64, page, size int) ([]Relation, error) // 关注者列表
	Detail(c context.Context, follower, followee int64) (Relation, error)              // 关注状态
	FollowerCnt(c context.Context, userID int64) (int64, error)                        // 粉丝数
	FolloweeCnt(c context.Context, userID int64) (int64, error)                        // 关注数
}
