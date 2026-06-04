package middleware

import (
	"BlockCertify/internal/services/CacheService"
	"BlockCertify/internal/services/LogService"

	"gorm.io/gorm"
)

type Middlewares struct {
	DB             *gorm.DB
	cacheService   *CacheService.CacheService
	logBulkService *LogService.LogBulkService
}
