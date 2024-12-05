package response

import (
	"github.com/gin-gonic/gin"
)

func Error(c *gin.Context, status int, err error, message string) {
	errorRes := message

	c.JSON(status, Response{Error: errorRes})
}
