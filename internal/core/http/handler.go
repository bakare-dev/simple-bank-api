package routes

import (
	"github.com/bakare-dev/simple-bank-api/internal/core/service"
	"github.com/gin-gonic/gin"
)

type CoreHandler struct {
	CoreService *service.CoreService
}

func NewCoreHandler(CoreService *service.CoreService) *CoreHandler {
	return &CoreHandler{CoreService: CoreService}
}

func (h *CoreHandler) HandleCreateUser(ctx *gin.Context) {

}
