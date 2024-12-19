package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"hr-system/internal/cache"
	common_errors "hr-system/internal/common/errors"
	"hr-system/internal/leaves/domain"
)

type LeaveCache interface {
	SetLeaveToCache(ctx context.Context, leave *domain.Leave) error
	GetLeaveFromCache(ctx context.Context, id int) (domain.Leave, error)
	SetLeavesToCache(ctx context.Context, query domain.LeavesQuery, leaves []domain.Leave) error
	GetLeavesFromCache(ctx context.Context, query domain.LeavesQuery) ([]domain.Leave, error)
	DelLeaveFromCache(ctx context.Context, id int) error
	DelLeavesFromCache(ctx context.Context, query domain.LeavesQuery) error
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

func (c *leaveCache) genCacheKeyByID(id int) string {
	return fmt.Sprintf("%s_leave_%d", c.prefix, id)
}

func (c *leaveCache) SetLeaveToCache(ctx context.Context, leave *domain.Leave) error {
	cacheKey := c.genCacheKeyByID(leave.ID)
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

func (c *leaveCache) GetLeaveFromCache(ctx context.Context, id int) (domain.Leave, error) {
	cacheKey := c.genCacheKeyByID(id)
	data, err := c.cache.Get(ctx, cacheKey)
	if errors.Is(err, redis.Nil) {
		return domain.Leave{}, common_errors.ErrResourceNotFound
	} else if err != nil {
		return domain.Leave{}, err
	}

	var leave domain.Leave
	err = json.Unmarshal([]byte(data), &leave)
	if err != nil {
		return domain.Leave{}, err
	}

	return leave, nil
}

func (c *leaveCache) SetLeavesToCache(ctx context.Context, query domain.LeavesQuery, leaves []domain.Leave) error {
	cacheKey := c.genCacheKey(query)
	data, err := json.Marshal(leaves)
	if err != nil {
		return err
	}

	err = c.cache.Set(ctx, cacheKey, string(data), 1*time.Hour)
	if err != nil {
		return err
	}
	return nil
}

func (c *leaveCache) genCacheKey(query domain.LeavesQuery) string {
	if query.EmployeeID != nil {
		return fmt.Sprintf("%s_employee_%d", c.prefix, *query.EmployeeID)
	} else if query.CurrentReviewerID != nil {
		return fmt.Sprintf("%s_reviewer_%d", c.prefix, *query.CurrentReviewerID)
	} else {
		// won't happen now
		return fmt.Sprintf("%s_all", c.prefix)
	}
}

func (c *leaveCache) GetLeavesFromCache(ctx context.Context, query domain.LeavesQuery) ([]domain.Leave, error) {
	cacheKey := c.genCacheKey(query)
	data, err := c.cache.Get(ctx, cacheKey)
	if errors.Is(err, redis.Nil) {
		return nil, common_errors.ErrResourceNotFound
	} else if err != nil {
		return nil, err
	}
	if data == "" {
		return nil, common_errors.ErrResourceNotFound
	}

	var leaves []domain.Leave
	err = json.Unmarshal([]byte(data), &leaves)
	if err != nil {
		return nil, err
	}

	return leaves, nil
}

func (c *leaveCache) DelLeavesFromCache(ctx context.Context, query domain.LeavesQuery) error {
	cacheKey := c.genCacheKey(query)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

func (c *leaveCache) DelLeaveFromCache(ctx context.Context, id int) error {
	cacheKey := c.genCacheKeyByID(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
