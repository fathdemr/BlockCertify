package CacheService

import (
	"context"
)

func (s *CacheService) Get(key string) (string, error) {
	return s.redi.Get(context.Background(), key).Result()
}
