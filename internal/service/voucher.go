package service

import (
	"context"
	"fmt"
	"github.com/jinzhu/copier"
	v1 "go-dianping/api/v1"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/model"
)

type VoucherService interface {
	AddSeckillVoucher(ctx context.Context, req *model.Voucher) error
	AddVoucher(ctx context.Context, req *model.Voucher) error
	QueryVoucherOfShop(ctx context.Context, req *v1.QueryVoucherOfShopReq) ([]*model.Voucher, error)
}

func NewVoucherService(
	service *Service,
) VoucherService {
	return &voucherService{
		Service: service,
	}
}

type voucherService struct {
	*Service
}

func (s *voucherService) AddSeckillVoucher(ctx context.Context, voucher *model.Voucher) error {
	// 保存优惠券
	if err := s.query.Voucher.Save(voucher); err != nil {
		return err
	}
	// 保存秒杀信息
	var seckillVoucher model.SeckillVoucher
	seckillVoucher.VoucherID = voucher.ID
	if err := copier.Copy(&seckillVoucher, &voucher); err != nil {
		return err
	}
	if err := s.query.SeckillVoucher.Save(&seckillVoucher); err != nil {
		return err
	}
	// 保存秒杀库存到 redis 中
	key := fmt.Sprintf("%s%d", constants.RedisSeckillStockKey, voucher.ID)
	return s.rdb.Set(ctx, key, voucher.Stock, 0).Err()
}

func (s *voucherService) AddVoucher(_ context.Context, voucher *model.Voucher) error {
	return s.query.Voucher.Save(voucher)
}

func (s *voucherService) QueryVoucherOfShop(_ context.Context, req *v1.QueryVoucherOfShopReq) ([]*model.Voucher, error) {
	return s.query.Voucher.Where(s.query.Voucher.ShopID.Eq(req.ShopId)).Find()
}
