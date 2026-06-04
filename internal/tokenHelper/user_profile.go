package tokenHelper

import (
	"github.com/gin-gonic/gin"
)

func GetActor(c *gin.Context) *Actor {
	var actor Actor
	if m, isExist := c.Get("actor"); isExist {
		actor = m.(Actor)
	}
	return &actor
}
