package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"hr-system/internal/cache"
	"hr-system/internal/leaves/domain"
)

type LeaveCache interface {
	SetLeaveToCache(ctx context.Context, leave *domain.Leave) error
	GetLeaveFromCache(ctx context.Context, id int) (*domain.Leave, error)
	SetLeavesToCache(ctx context.Context, leaves []domain.Leave, page, pageSize int) error
	GetLeavesFromCache(ctx context.Context, page, pageSize int) ([]domain.Leave, error)
}

type leaveCache struct {
	cache  *cache.Cache
	prefix string
}

func NewLeaveCache(cache *cache.Cache, prefix string) LeaveCache {
	return &leaveCache{
		cache:  cache,
		prefix: prefix,
	}
}

func (c *leaveCache) genCacheKey(id int) string {
	return fmt.Sprintf("%s_leave_%d", c.prefix, id)
}

func (c *leaveCache) SetLeaveToCache(ctx context.Context, leave *domain.Leave) error {
	cacheKey := c.genCacheKey(leave.ID)
	data, err := json.Marshal(leave)
	if err != nil {
		return err
	}

	err = c.cache.Set(ctx, cacheKey, string(data), 1*time.Hour)
	if err != nil {
		return err
	}
	return nil
}

func (c *leaveCache) GetLeaveFromCache(ctx context.Context, id int) (*domain.Leave, error) {
	return nil, nil
	//cacheKey := c.genCacheKey(id)
	//data, err := c.cache.Get(ctx, cacheKey).Result()
	//if err == redis.Nil {
	//	return nil, common.ErrResourceNotFound
	//} else if err != nil {
	//	return nil, err
	//}
	//
	//var leave domain.Leave
	//err = json.Unmarshal([]byte(data), &leave)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return &leave, nil
}

func (c *leaveCache) SetLeavesToCache(ctx context.Context, leaves []domain.Leave, page, pageSize int) error {
	return nil
	//cacheKey := fmt.Sprintf("%s_leaves_page_%d_%d", c.prefix, page, pageSize)
	//data, err := json.Marshal(leaves)
	//if err != nil {
	//	return err
	//}
	//
	//err = c.client.Set(ctx, cacheKey, data, 24*time.Hour).Err()
	//if err != nil {
	//	return err
	//}
	//return nil
}

func (c *leaveCache) GetLeavesFromCache(ctx context.Context, page, pageSize int) ([]domain.Leave, error) {
	return nil, nil
	//cacheKey := fmt.Sprintf("%s_leaves_page_%d_%d", c.prefix, page, pageSize)
	//data, err := c.client.Get(ctx, cacheKey).Result()
	//if err == redis.Nil {
	//	return nil, common.ErrResourceNotFound
	//} else if err != nil {
	//	return nil, err
	//}
	//
	//var leaves []domain.Leave
	//err = json.Unmarshal([]byte(data), &leaves)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return leaves, nil
}
