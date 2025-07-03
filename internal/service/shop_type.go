package service

import (
	"context"
	"encoding/json"
	"github.com/samber/lo"
	"go-dianping/internal/base/constants"
	"go-dianping/internal/model"
)

type ShopTypeService interface {
	QueryTypeList(ctx context.Context) ([]*model.ShopType, error)
}

type shopTypeService struct {
	*Service
}

func NewShopTypeService(
	service *Service,
) ShopTypeService {
	return &shopTypeService{
		Service: service,
	}
}

func (s *shopTypeService) QueryTypeList(ctx context.Context) ([]*model.ShopType, error) {
	// 使用 List 结构

	exist, err := s.rdb.Exists(ctx, constants.RedisCacheShopTypeKey).Result()
	if err != nil {
		return nil, err
	}
	if exist == 1 {
		strSlices, err := s.rdb.LRange(ctx, constants.RedisCacheShopTypeKey, 0, -1).Result()
		if err != nil {
			return nil, err
		}
		return lo.Map(strSlices, func(el string, idx int) *model.ShopType {
			var item model.ShopType
			if err := json.Unmarshal([]byte(el), &item); err != nil {
				return &item
			}
			return &item
		}), nil
	}

	list, err := s.query.ShopType.Order(s.query.ShopType.Sort.Asc()).Find()
	if err != nil {
		return nil, err
	}

	lo.ForEach(list, func(el *model.ShopType, idx int) {
		bytes, closureErr := json.Marshal(el)
		if closureErr != nil {
			err = closureErr
		}
		if closureErr := s.rdb.RPush(ctx, constants.RedisCacheShopTypeKey, string(bytes)).Err(); closureErr != nil {
			err = closureErr
		}

	})

	if err := s.rdb.Expire(ctx, constants.RedisCacheShopTypeKey, constants.RedisCacheShopTypeTTL).Err(); err != nil {
		return nil, err
	}

	return list, nil
}
