package response

import (
	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, err error, message string) {
	errorRes := map[string]interface{}{
		"message": message,
	}

	c.JSON(status, Response{Error: errorRes})
}
