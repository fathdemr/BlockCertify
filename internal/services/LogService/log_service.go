package LogService

import (
	"BlockCertify/internal/models"
	"sync"
	"time"

	"gorm.io/gorm"
)

type LogBulkService struct {
	logsMap  map[string][]models.Log
	syncerMu sync.Mutex
	logS     *LogService
}

func NewLogBulkService(db *gorm.DB) *LogBulkService {
	return &LogBulkService{
		logS: New(db),
	}
}

type LogService struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *LogService {
	return &LogService{
		DB: db,
	}
}

func (l *LogBulkService) AppendAndBulkInsert(log models.Log) error {
	key := time.Now().Format("2006.01.02_15")
	l.syncerMu.Lock()
	if l.logsMap == nil {
		l.logsMap = make(map[string][]models.Log)
	}
	if _, ok := l.logsMap[key]; ok {
		l.logsMap[key] = append(l.logsMap[key], log)
	} else {
		l.logsMap[key] = []models.Log{log}
	}
	l.syncerMu.Unlock()
	return nil
}
