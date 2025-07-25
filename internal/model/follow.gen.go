// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameFollow = "tb_follow"

// Follow mapped from table <tb_follow>
type Follow struct {
	ID           int64      `gorm:"column:id;type:bigint(20);primaryKey;autoIncrement:true;comment:主键" json:"id"`                                // 主键
	UserID       uint64     `gorm:"column:user_id;type:bigint(20) unsigned;not null;comment:用户id" json:"userId"`                                 // 用户id
	FollowUserID uint64     `gorm:"column:follow_user_id;type:bigint(20) unsigned;not null;comment:关联的用户id" json:"followUserId"`                 // 关联的用户id
	CreateTime   *time.Time `gorm:"column:create_time;type:timestamp;not null;default:current_timestamp();autoCreateTime;comment:创建时间" json:"-"` // 创建时间
}

// TableName Follow's table name
func (*Follow) TableName() string {
	return TableNameFollow
}
