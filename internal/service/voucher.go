package service

import (
	"context"
	"github.com/jinzhu/copier"
	v1 "go-dianping/api/v1"
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

func (s *voucherService) AddSeckillVoucher(_ context.Context, req *model.Voucher) error {
	// 保存优惠券
	if err := s.query.Voucher.Save(req); err != nil {
		return err
	}
	// 保存秒杀信息
	var seckillVoucher model.SeckillVoucher
	seckillVoucher.VoucherID = req.ID
	if err := copier.Copy(&seckillVoucher, &req); err != nil {
		return err
	}
	return s.query.SeckillVoucher.Save(&seckillVoucher)
}

func (s *voucherService) AddVoucher(_ context.Context, req *model.Voucher) error {
	return s.query.Voucher.Save(req)
}

func (s *voucherService) QueryVoucherOfShop(_ context.Context, req *v1.QueryVoucherOfShopReq) ([]*model.Voucher, error) {
	return s.query.Voucher.Where(s.query.Voucher.ShopID.Eq(req.ShopId)).Find()
}
