package service

import (
	"context"
	"go-dianping/api/v1"
	"go-dianping/internal/base/redis_worker"
	"go-dianping/internal/base/user_holder"
	"go-dianping/internal/model"
	"go-dianping/internal/query"
	"time"
)

type VoucherOrderService interface {
	SeckillVoucher(ctx context.Context, req *v1.SeckillVoucherReq) (int64, error)
}

func NewVoucherOrderService(
	service *Service,
	redisWorker redis_worker.RedisWorker,
) VoucherOrderService {
	return &voucherOrderService{
		Service:     service,
		redisWorker: redisWorker,
	}
}

type voucherOrderService struct {
	*Service
	redisWorker redis_worker.RedisWorker
}

func (s *voucherOrderService) SeckillVoucher(ctx context.Context, req *v1.SeckillVoucherReq) (int64, error) {
	//	1. 查询优惠券
	voucher, err := s.query.SeckillVoucher.Where(s.query.SeckillVoucher.VoucherID.Eq(req.ID)).First()
	if err != nil {
		return 0, err
	}
	//	2. 判断秒杀是否开始
	if voucher.BeginTime.After(time.Now()) {
		return 0, v1.ErrSeckillNotStart
	}
	//	3. 判断秒杀是否已经结束
	if voucher.EndTime.Before(time.Now()) {
		return 0, v1.ErrSeckillIsEnd
	}
	//	4. 判断库存是否充足
	if voucher.Stock < 1 {
		return 0, v1.ErrInsufficientStock
	}

	return s.createVoucherOrder(ctx, err, voucher)
}

func (s *voucherOrderService) createVoucherOrder(ctx context.Context, err error, voucher *model.SeckillVoucher) (int64, error) {
	// 5. 一人一单
	userId := user_holder.GetUser(ctx).ID
	// 5.1. 查询订单
	count, err := s.query.VoucherOrder.Where(s.query.VoucherOrder.UserID.Eq(*userId)).
		Where(s.query.VoucherOrder.VoucherID.Eq(voucher.VoucherID)).Count()
	if err != nil {
		return 0, err
	}
	// 5.2. 判断是否存在
	if count > 0 {
		// 用户已经购买过了
		return 0, v1.ErrAlreadySeckill
	}

	//	6. 扣减库存，返回订单 id
	var orderId int64
	return orderId, s.query.Transaction(func(tx *query.Query) error {
		info, err := s.query.WithContext(ctx).SeckillVoucher.
			Where(s.query.SeckillVoucher.VoucherID.Eq(voucher.VoucherID)).
			Where(s.query.SeckillVoucher.Stock.Gt(0)).
			Update(s.query.SeckillVoucher.Stock, s.query.SeckillVoucher.Stock.Sub(1))
		if err != nil {
			return err
		}
		if info.Error != nil {
			// 扣减失败
			return v1.ErrInsufficientStock
		}
		//	7. 创建订单
		var voucherOrder model.VoucherOrder
		// 7.1. 订单 id
		orderId, err = s.redisWorker.NextId(ctx, "order")
		if err != nil {
			return err
		}
		voucherOrder.ID = orderId
		// 7.2. 用户 id
		voucherOrder.UserID = *userId
		// 7.3. 代金券 id
		voucherOrder.VoucherID = voucher.VoucherID
		if err := s.query.VoucherOrder.Save(&voucherOrder); err != nil {
			return err
		}

		return nil
	})
}
