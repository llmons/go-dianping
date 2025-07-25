package service

import (
	"context"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"go-dianping/api/v1"
	"go-dianping/internal/base/redis_worker"
	"go-dianping/internal/base/user_holder"
	"go-dianping/internal/model"
	"go-dianping/internal/query"
	"go.uber.org/zap"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const streamName = "stream.orders"

type VoucherOrderService interface {
	SeckillVoucher(ctx context.Context, req *v1.SeckillVoucherReq) (int64, error)
}

type voucherOrderService struct {
	*Service
	redisWorker   redis_worker.RedisWorker
	seckillScript *redis.Script
	orderTasks    chan *model.VoucherOrder
}

func NewVoucherOrderService(
	service *Service,
	redisWorker redis_worker.RedisWorker,
) VoucherOrderService {
	// 加载 lua 脚本
	workDir, err := os.Getwd()
	if err != nil {
		return nil
	}
	luaPath := filepath.Join(workDir, "internal", "scripts", "seckill.lua")
	bytes, err := os.ReadFile(luaPath)
	if err != nil {
		return nil
	}
	seckillScript := redis.NewScript(string(bytes))

	orderTasks := make(chan *model.VoucherOrder, 1024*1024)
	srv := &voucherOrderService{
		Service:       service,
		redisWorker:   redisWorker,
		seckillScript: seckillScript,
		orderTasks:    orderTasks,
	}
	go func() {
		for {
			// 1. 获取消息队列中的订单信息
			ctx := context.Background()
			result, err := srv.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    "g1",
				Consumer: "c1",
				Count:    1,
				Block:    time.Second * 2,
				Streams:  []string{streamName, ">"},
			}).Result()
			if err != nil {
				continue
			}
			// 2. 判断消息获取是否成功
			if len(result) == 0 {
				//	如果获取失败，说明没有消息，继续下一次循环
				continue
			}
			// 3. 解析消息中的订单信息
			messages := result[0].Messages
			values := messages[0].Values
			id, err := strconv.Atoi(values["id"].(string))
			if err != nil {
				continue
			}
			userID, err := strconv.Atoi(values["userId"].(string))
			if err != nil {
				continue
			}
			voucherID, err := strconv.Atoi(values["voucherId"].(string))
			if err != nil {
				continue
			}
			order := model.VoucherOrder{
				ID:        int64(id),
				UserID:    uint64(userID),
				VoucherID: uint64(voucherID),
			}
			// 4. 如果获取成功，可以下单
			if err := srv.handleVoucherOrder(&order); err != nil {
				srv.logger.Error("处理订单异常", zap.Error(err))
				srv.handlePendingList()
			}
			// 5. ACK 消息
			if err := srv.rdb.XAck(ctx, streamName, "g1", messages[0].ID).Err(); err != nil {
				continue
			}
		}
	}()

	//orderTasks := make(chan *model.VoucherOrder, 1024*1024)
	//srv := &voucherOrderService{
	//	Service:       service,
	//	redisWorker:   redisWorker,
	//	seckillScript: seckillScript,
	//	orderTasks:    orderTasks,
	//}
	//go func() {
	//	for {
	//		// 1. 获取 chan 中的订单信息
	//		order := <-srv.orderTasks
	//		// 2. 创建订单
	//		if err := srv.handleVoucherOrder(order); err != nil {
	//			srv.logger.Error("处理订单异常", zap.Error(err))
	//		}
	//	}
	//}()

	return srv
}

func (s *voucherOrderService) handleVoucherOrder(order *model.VoucherOrder) (err error) {
	// 1. 获取用户
	userId := order.UserID
	// 2. 创建锁
	lockName := fmt.Sprintf("lock:order:%d", userId)
	lock := s.rs.NewMutex(lockName)
	// 获取锁
	err = lock.TryLock()
	// 判断是否获取锁成功
	if err != nil {
		//	获取锁失败，返回错误或重试
		s.logger.Error("不允许重复下单")
		return err
	}
	// 结束时释放锁
	defer func(lock *redsync.Mutex) {
		_, closureErr := lock.Unlock()
		if closureErr != nil {
			err = closureErr
		}
	}(lock)
	return s.createVoucherOrder(order)
}

func (s *voucherOrderService) handlePendingList() {
	for {
		// 1. 获取 pending list 中的订单信息
		ctx := context.Background()
		result, err := s.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    "g1",
			Consumer: "c1",
			Count:    1,
			Block:    time.Second * 2,
			Streams:  []string{streamName, "0"},
		}).Result()
		if err != nil {
			continue
		}
		// 2. 判断消息获取是否成功
		if len(result) == 0 {
			//	如果获取失败，说明 pending list 没有消息，结束循环
			break
		}
		// 3. 解析消息中的订单信息
		messages := result[0].Messages
		values := messages[0].Values
		id, err := strconv.Atoi(values["id"].(string))
		if err != nil {
			continue
		}
		userID, err := strconv.Atoi(values["userId"].(string))
		if err != nil {
			continue
		}
		voucherID, err := strconv.Atoi(values["voucherId"].(string))
		if err != nil {
			continue
		}
		order := model.VoucherOrder{
			ID:        int64(id),
			UserID:    uint64(userID),
			VoucherID: uint64(voucherID),
		}
		// 4. 如果获取成功，可以下单
		if err := s.handleVoucherOrder(&order); err != nil {
			s.logger.Error("处理订单异常", zap.Error(err))
			time.Sleep(time.Millisecond * 20)
			continue
		}
		// 5. ACK 消息
		if err := s.rdb.XAck(ctx, streamName, "g1", messages[0].ID).Err(); err != nil {
			continue
		}
	}
}

func (s *voucherOrderService) SeckillVoucher(ctx context.Context, req *v1.SeckillVoucherReq) (voucherId int64, err error) {
	// 获取用户 id
	userId := user_holder.GetUser(ctx).ID
	// 生成订单 id
	orderId, err := s.redisWorker.NextId(ctx, "order")
	if err != nil {
		return 0, err
	}
	// 1. 执行 lua 脚本
	result, err := s.seckillScript.Run(ctx, s.rdb, []string{}, req.VoucherID, *userId, orderId).Result()
	if err != nil {
		return 0, err
	}
	// 2. 判断结果是否为 0
	if result.(int64) != 0 {
		// 2.1. 不为 0，代表没有购买资格
		if result.(int64) == 1 {
			return 0, v1.ErrInsufficientStock
		} else {
			return 0, v1.ErrNotAllowDoubleBuy
		}
	}

	// 3. 返回订单 id
	return orderId, nil
}

//func (s *voucherOrderService) SeckillVoucher(ctx context.Context, req *v1.SeckillVoucherReq) (voucherId int64, err error) {
//	// 获取用户
//	userId := user_holder.GetUser(ctx).FollowUserId
//	1. 执行 lua 脚本
//	result, err := s.seckillScript.Run(ctx, s.rdb, []string{}, req.VoucherID, *userId).Result()
//	if err != nil {
//		return 0, err
//	}
//	// 2. 判断结果是否为 0
//	if result.(int64) != 0 {
//		// 2.1. 不为 0，代表没有购买资格
//		if result.(int64) == 1 {
//			return 0, v1.ErrInsufficientStock
//		} else {
//			return 0, v1.ErrNotAllowDoubleBuy
//		}
//	}
//	// 2.2. 为 0，有购买资格，把下单信息保存到阻塞队列
//	var voucherOrder model.VoucherOrder
//	// 2.3. 订单 id
//	orderId, err = s.redisWorker.NextId(ctx, "order")
//	if err != nil {
//		return 0, err
//	}
//	voucherOrder.FollowUserId = orderId
//	// 2.4. 用户 id
//	voucherOrder.UserID = *userId
//	// 2.5. 代金券 id
//	voucherOrder.VoucherID = req.VoucherID
//	// 2.6. 放入阻塞队列
//	s.orderTasks <- &voucherOrder
//
//	// 3. 返回订单 id
//	return orderId, nil
//}

//func (s *voucherOrderService) SeckillVoucher(ctx context.Context, req *v1.SeckillVoucherReq) (voucherId int64, err error) {
//	//	1. 查询优惠券
//	voucher, err := s.query.SeckillVoucher.Where(s.query.SeckillVoucher.VoucherID.Eq(req.VoucherID)).First()
//	if err != nil {
//		return 0, err
//	}
//	//	2. 判断秒杀是否开始
//	if voucher.BeginTime.After(time.Now()) {
//		return 0, v1.ErrSeckillNotStart
//	}
//	//	3. 判断秒杀是否已经结束
//	if voucher.EndTime.Before(time.Now()) {
//		return 0, v1.ErrSeckillIsEnd
//	}
//	//	4. 判断库存是否充足
//	if voucher.Stock < 1 {
//		// 库存不足
//		return 0, v1.ErrInsufficientStock
//	}
//
//	userId := user_holder.GetUser(ctx).VoucherID
//	lockName := fmt.Sprintf("order:%d", *userId)
//
//	// ==================== CUSTOM LOCK ====================
//	// 创建锁
//	//lock := redis_lock.NewSimpleRedisLock(lockName, s.rdb)
//	// 获取锁
//	//isLock, err := lock.TryLock(ctx, time.Second*5)
//	//if err != nil {
//	//	return 0, err
//	//}
//	// 判断是否获取锁成功
//	//if !isLock {
//	//	//	获取锁失败，返回错误或重试
//	//	return 0, v1.ErrNotAllowDoubleBuy
//	//}
//	// 结束时释放锁
//	//defer func(lock redis_lock.ILock, ctx context.Context) {
//	//	closureErr := lock.Unlock(ctx)
//	//	if closureErr != nil {
//	//		err = closureErr
//	//	}
//	//}(lock, ctx)
//
//	// ==================== RedSync Lib ====================
//	// 创建锁
//	lock := s.rs.NewMutex(lockName)
//	// 获取锁
//	err = lock.TryLock()
//	// 判断是否获取锁成功
//	if err != nil {
//		//	获取锁失败，返回错误或重试
//		return 0, v1.ErrNotAllowDoubleBuy
//	}
//	// 结束时释放锁
//	defer func(lock *redsync.Mutex) {
//		_, closureErr := lock.Unlock()
//		if closureErr != nil {
//			err = closureErr
//		}
//	}(lock)
//	return s.createVoucherOrder(ctx, err, voucher)
//}

func (s *voucherOrderService) createVoucherOrder(voucherOrder *model.VoucherOrder) error {
	return s.query.Transaction(func(tx *query.Query) error {
		//	5. 一人一单
		userId := voucherOrder.UserID

		//	5.1. 查询订单
		count, err := s.query.VoucherOrder.Where(s.query.VoucherOrder.UserID.Eq(userId)).
			Where(s.query.VoucherOrder.VoucherID.Eq(voucherOrder.VoucherID)).Count()
		if err != nil {
			return err
		}
		//	5.2.判断是否存在
		if count > 0 {
			//	用户已经购买过了
			s.logger.Error("用户已经购买过一次")
			return v1.ErrAlreadySeckill
		}
		//	6. 扣减库存，返回订单 id
		info, err := s.query.SeckillVoucher.
			Where(s.query.SeckillVoucher.VoucherID.Eq(voucherOrder.VoucherID)).
			Where(s.query.SeckillVoucher.Stock.Gt(0)).
			Update(s.query.SeckillVoucher.Stock, s.query.SeckillVoucher.Stock.Sub(1))
		if err != nil {
			return err
		}
		if info.Error != nil || info.RowsAffected == 0 { // info.RowsAffected == 0 即库存为 0，没有更新（扣减失败），但是也没 error
			// 扣减失败
			s.logger.Error("库存不足")
			return v1.ErrInsufficientStock
		}
		//	7. 创建订单
		if err := s.query.VoucherOrder.Save(voucherOrder); err != nil {
			return err
		}

		return nil
	})
}
