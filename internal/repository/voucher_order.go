package repository

import (
	"context"
	"go-dianping/internal/model"
)

type VoucherOrderRepository interface {
	Save(context.Context, *model.VoucherOrder) error
}

func NewVoucherOrderRepository(
	repository *Repository,
) VoucherOrderRepository {
	return &voucherOrderRepository{
		Repository: repository,
	}
}

type voucherOrderRepository struct {
	*Repository
}

func (r *voucherOrderRepository) Save(ctx context.Context, order *model.VoucherOrder) error {
	return r.query.WithContext(ctx).VoucherOrder.Save(order)
}
