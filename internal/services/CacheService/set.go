package CacheService

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func (s *CacheService) SetWithExpireDuration(key string, value any, duration time.Duration) error {
	return s.redi.Set(context.Background(), key, value, duration).Err()
}

func (s *CacheService) UpdateKeyValue(key string, value interface{}) error {
	return s.redi.SetArgs(context.Background(), key, value, redis.SetArgs{
		KeepTTL: true,
	}).Err()
}

func (s *CacheService) ExtendTTL(key string, value time.Duration) error {

	ctx := context.Background()

	ttl, err := s.redi.TTL(ctx, key).Result()
	if err != nil {
		return err
	}

	// TTL DEĞERİ -2 -> Key yok -1 -> Süresiz Key >0 -> Süreli Key
	if ttl == -1*time.Second {
		return fmt.Errorf("süresiz key")
	}
	if ttl == -2*time.Second {
		return fmt.Errorf("key yok")
	}

	newTTL := value + ttl

	return s.redi.Expire(ctx, key, newTTL).Err()
}
