package CacheService

import (
	"context"
	"strings"
)

func (s *CacheService) Delete(key string) error {
	return s.redi.Del(context.Background(), key).Err()
}

func (s *CacheService) DeleteMultipleKeys(key ...string) error {
	return s.redi.Del(context.Background(), key...).Err()
}

func (s *CacheService) DeleteAll() error {
	return s.redi.FlushAll(context.Background()).Err()
}

func (s *CacheService) DeleteStartWithKeys(prefix string) {
	if prefix == "" {
		return
	}
	keys, _ := s.GetAllKeys()
	for _, key := range keys {
		if strings.HasPrefix(key, prefix) {
			s.redi.Del(context.Background(), key)
		}
	}
}
