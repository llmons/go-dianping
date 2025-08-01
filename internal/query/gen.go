// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:             db,
		Blog:           newBlog(db, opts...),
		BlogComments:   newBlogComments(db, opts...),
		Follow:         newFollow(db, opts...),
		SeckillVoucher: newSeckillVoucher(db, opts...),
		Shop:           newShop(db, opts...),
		ShopType:       newShopType(db, opts...),
		Sign:           newSign(db, opts...),
		User:           newUser(db, opts...),
		UserInfo:       newUserInfo(db, opts...),
		Voucher:        newVoucher(db, opts...),
		VoucherOrder:   newVoucherOrder(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	Blog           blog
	BlogComments   blogComments
	Follow         follow
	SeckillVoucher seckillVoucher
	Shop           shop
	ShopType       shopType
	Sign           sign
	User           user
	UserInfo       userInfo
	Voucher        voucher
	VoucherOrder   voucherOrder
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:             db,
		Blog:           q.Blog.clone(db),
		BlogComments:   q.BlogComments.clone(db),
		Follow:         q.Follow.clone(db),
		SeckillVoucher: q.SeckillVoucher.clone(db),
		Shop:           q.Shop.clone(db),
		ShopType:       q.ShopType.clone(db),
		Sign:           q.Sign.clone(db),
		User:           q.User.clone(db),
		UserInfo:       q.UserInfo.clone(db),
		Voucher:        q.Voucher.clone(db),
		VoucherOrder:   q.VoucherOrder.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:             db,
		Blog:           q.Blog.replaceDB(db),
		BlogComments:   q.BlogComments.replaceDB(db),
		Follow:         q.Follow.replaceDB(db),
		SeckillVoucher: q.SeckillVoucher.replaceDB(db),
		Shop:           q.Shop.replaceDB(db),
		ShopType:       q.ShopType.replaceDB(db),
		Sign:           q.Sign.replaceDB(db),
		User:           q.User.replaceDB(db),
		UserInfo:       q.UserInfo.replaceDB(db),
		Voucher:        q.Voucher.replaceDB(db),
		VoucherOrder:   q.VoucherOrder.replaceDB(db),
	}
}

type queryCtx struct {
	Blog           IBlogDo
	BlogComments   IBlogCommentsDo
	Follow         IFollowDo
	SeckillVoucher ISeckillVoucherDo
	Shop           IShopDo
	ShopType       IShopTypeDo
	Sign           ISignDo
	User           IUserDo
	UserInfo       IUserInfoDo
	Voucher        IVoucherDo
	VoucherOrder   IVoucherOrderDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		Blog:           q.Blog.WithContext(ctx),
		BlogComments:   q.BlogComments.WithContext(ctx),
		Follow:         q.Follow.WithContext(ctx),
		SeckillVoucher: q.SeckillVoucher.WithContext(ctx),
		Shop:           q.Shop.WithContext(ctx),
		ShopType:       q.ShopType.WithContext(ctx),
		Sign:           q.Sign.WithContext(ctx),
		User:           q.User.WithContext(ctx),
		UserInfo:       q.UserInfo.WithContext(ctx),
		Voucher:        q.Voucher.WithContext(ctx),
		VoucherOrder:   q.VoucherOrder.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}
