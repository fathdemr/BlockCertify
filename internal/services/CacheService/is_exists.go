package CacheService

import "context"

func (s *CacheService) IsExists(key string) (bool, error) {
	val, err := s.redi.Exists(context.Background(), key).Result()
	if err != nil {
		return false, err
	}
	return val > 0, err
}
