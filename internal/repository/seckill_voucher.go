package repository

import (
	"context"
	"go-dianping/internal/entity"
)

type SeckillVoucherRepository interface {
	Save(ctx context.Context, seckillVoucher *entity.SeckillVoucher) error
}

func NewSeckillVoucherRepository(
	repository *Repository,
) SeckillVoucherRepository {
	return &seckillVoucherRepository{
		Repository: repository,
	}
}

type seckillVoucherRepository struct {
	*Repository
}

func (r *seckillVoucherRepository) Save(ctx context.Context, seckillVoucher *entity.SeckillVoucher) error {
	return r.query.WithContext(ctx).SeckillVoucher.Save(seckillVoucher)
}
