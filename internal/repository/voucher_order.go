package repository

import (
	"context"
	"go-dianping/internal/entity"
)

type VoucherOrderRepository interface {
	Save(context.Context, *entity.VoucherOrder) error
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

func (r *voucherOrderRepository) Save(ctx context.Context, order *entity.VoucherOrder) error {
	return r.query.WithContext(ctx).VoucherOrder.Save(order)
}
