package CacheService

import (
	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	redi *redis.Client
}

func New(redi *redis.Client) *CacheService {
	return &CacheService{
		redi: redi,
	}
}
