package webserver

import (
	"errors"
	"github.com/RollingBalls/rollingballs-server/engine"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CheckAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token := c.Request.Header.Get("X-Auth-Token"); token != "" {
			c.Next()
		} else {
			c.Fail(http.StatusUnauthorized, errors.New("Unauthorized"))
		}
	}
}

func RollingBalls(engine *engine.Engine) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("engine", engine)

		c.Next()
	}
}
