package middleware

import (
	"BlockCertify/internal/models"
	"BlockCertify/internal/tokenHelper"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func (m *Middlewares) Logger(c *gin.Context) {
	startTime := time.Now()
	const responseBodyKey = "x-log-response-body"
	const requestBodyKey = "x-log-request-body"
	c.Next()
	if c.Request.Method == "OPTIONS" {
		return
	}

	// Log işlemi olmamaması durumunda direkt çık
	var exist bool
	var logEnableAny any
	logEnableAny, exist = c.Get("x-log-enable")
	if exist && logEnableAny != nil {
		disableCheck, ok := logEnableAny.(bool)
		if ok && !disableCheck {
			return
		}
	}

	actor := tokenHelper.GetActor(c)

	var log models.Log
	log.IpAddress = c.ClientIP()
	log.ID = actor.ID
	log.HTTPMethod = c.Request.Method
	log.URLPath = c.Request.URL.Path
	log.ResponseStatus = c.Writer.Status()

	if reqBodyAny, isExist := c.Get(requestBodyKey); isExist {
		var reqBodyJSONBytes []byte
		reqBodyJSONBytes, err := json.Marshal(reqBodyAny)
		if err == nil {
			log.RequestBody = string(reqBodyJSONBytes)
		}
	}

	if resBodyAny, isExist := c.Get(responseBodyKey); isExist {
		var resBodyJSONBytes []byte
		resBodyJSONBytes, err := json.Marshal(resBodyAny)
		if err == nil {
			log.ResponseBody = string(resBodyJSONBytes)
		}
	}

	go func(l models.Log, startTime time.Time) {
		l.Latency = time.Since(startTime).Milliseconds()
		err := m.logBulkService.AppendAndBulkInsert(l)
		if err != nil {
			fmt.Printf("logService.Create returned error: %v\n", err)
		}
	}(log, startTime)
}
