package helper

import (
	"BlockCertify/internal/models"

	"github.com/gin-gonic/gin"
)

func SetSaveChangeLog(c *gin.Context,
	old, new any,
	action models.ChangeLogType, tableName string, relatedID uint64,
	ignoreFields []string) {
	c.Set("x-save-changelog", true)
	c.Set("x-changelog-table", tableName)
	c.Set("x-changelog-related-id", relatedID)
	c.Set("x-changelog-type", action)
	c.Set("x-changelog-old", old)
	c.Set("x-changelog-new", new)
	c.Set("x-changelog-ignored-fields", ignoreFields)
}

func SetLogResponse(c *gin.Context, response any) {
	c.Set("x-log-response-body", response)
}

func SetLogRequest(c *gin.Context, request any) {
	c.Set("x-log-request-body", request)
}

func SetLogEnable(c *gin.Context, enable bool) {
	c.Set("x-log-enable", enable)
}

func GinJSON(c *gin.Context,
	status int,
	response any,
	request any) {
	c.Set("x-log-response-body", response)
	c.Set("x-log-request-body", request)
	c.Set("x-log-enable", true)
	c.JSON(status, response)
}

func GinHTML(c *gin.Context,
	status int,
	name string,
	data any) {
	c.Set("x-log-response-body", data)
	c.Set("x-log-request-body", name)
	c.Set("x-log-enable", true)
	c.HTML(status, name, data)
}
