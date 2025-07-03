package repository

import (
	"context"
	"go-dianping/internal/entity"
	"gorm.io/gen"
)

type SeckillVoucherRepository interface {
	GetByID(ctx context.Context, id uint64) (*entity.SeckillVoucher, error)
	Save(ctx context.Context, seckillVoucher *entity.SeckillVoucher) error
	DecStock(ctx context.Context, id uint64) (gen.ResultInfo, error)
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

func (r *seckillVoucherRepository) GetByID(ctx context.Context, id uint64) (*entity.SeckillVoucher, error) {
	return r.query.WithContext(ctx).SeckillVoucher.Where(r.query.SeckillVoucher.VoucherID.Eq(id)).First()
}

func (r *seckillVoucherRepository) Save(ctx context.Context, seckillVoucher *entity.SeckillVoucher) error {
	return r.query.WithContext(ctx).SeckillVoucher.Save(seckillVoucher)
}

func (r *seckillVoucherRepository) DecStock(ctx context.Context, id uint64) (gen.ResultInfo, error) {
	return r.query.WithContext(ctx).SeckillVoucher.
		Where(r.query.SeckillVoucher.VoucherID.Eq(id)).
		Where(r.query.SeckillVoucher.Stock.Gt(0)).
		Update(r.query.SeckillVoucher.Stock, r.query.SeckillVoucher.Stock.Sub(1))
}
