package repository

import (
	"context"
	"go-dianping/internal/model"
)

type VoucherRepository interface {
	Save(ctx context.Context, voucher *model.Voucher) error
}

func NewVoucherRepository(
	repository *Repository,
) VoucherRepository {
	return &voucherRepository{
		Repository: repository,
	}
}

type voucherRepository struct {
	*Repository
}

func (r *voucherRepository) Save(ctx context.Context, voucher *model.Voucher) error {
	return r.query.WithContext(ctx).Voucher.Save(voucher)
}
