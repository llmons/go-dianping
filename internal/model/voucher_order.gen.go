// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameVoucherOrder = "tb_voucher_order"

// VoucherOrder mapped from table <tb_voucher_order>
type VoucherOrder struct {
	ID         int64      `gorm:"column:id;type:bigint(20);primaryKey;comment:主键" json:"id"`                                                                // 主键
	UserID     uint64     `gorm:"column:user_id;type:bigint(20) unsigned;not null;comment:下单的用户id" json:"userId"`                                           // 下单的用户id
	VoucherID  uint64     `gorm:"column:voucher_id;type:bigint(20) unsigned;not null;comment:购买的代金券id" json:"voucherId"`                                    // 购买的代金券id
	PayType    *uint8     `gorm:"column:pay_type;type:tinyint(1) unsigned;not null;default:1;comment:支付方式 1：余额支付；2：支付宝；3：微信" json:"payType"`                // 支付方式 1：余额支付；2：支付宝；3：微信
	Status     *uint8     `gorm:"column:status;type:tinyint(1) unsigned;not null;default:1;comment:订单状态，1：未支付；2：已支付；3：已核销；4：已取消；5：退款中；6：已退款" json:"status"` // 订单状态，1：未支付；2：已支付；3：已核销；4：已取消；5：退款中；6：已退款
	CreateTime *time.Time `gorm:"column:create_time;type:timestamp;not null;default:current_timestamp();autoCreateTime;comment:下单时间" json:"-"`              // 下单时间
	PayTime    *time.Time `gorm:"column:pay_time;type:timestamp;comment:支付时间" json:"payTime"`                                                               // 支付时间
	UseTime    *time.Time `gorm:"column:use_time;type:timestamp;comment:核销时间" json:"useTime"`                                                               // 核销时间
	RefundTime *time.Time `gorm:"column:refund_time;type:timestamp;comment:退款时间" json:"refundTime"`                                                         // 退款时间
	UpdateTime *time.Time `gorm:"column:update_time;type:timestamp;not null;default:current_timestamp();autoUpdateTime;comment:更新时间" json:"-"`              // 更新时间
}

// TableName VoucherOrder's table name
func (*VoucherOrder) TableName() string {
	return TableNameVoucherOrder
}
