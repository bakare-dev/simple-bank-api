package response

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
	Message interface{} `json:"message"`
}

func JSON(c *gin.Context, status int, data interface{}, message string) {
	c.JSON(status, Response{
		Data:    data,
		Message: message,
	})
}
