package routes

import (
	"net/http"

	"github.com/bakare-dev/simple-bank-api/internal/user/dto"
	"github.com/bakare-dev/simple-bank-api/internal/user/service"
	"github.com/bakare-dev/simple-bank-api/pkg/response"
	"github.com/gin-gonic/gin"
)

func handleCreateUser(ctx *gin.Context, userService *service.UserService) {
	var createUserDTO dto.CreateUserRequest

	if err := ctx.ShouldBindJSON(&createUserDTO); err != nil {
		response.Error(ctx, http.StatusBadRequest, err, err.Error())
		return
	}

	user := createUserDTO.ToModel()

	createdUser, status, err := userService.CreateUser(ctx.Request.Context(), user)
	if err != nil {
		response.Error(ctx, status, err, err.Error())
		return
	}

	modelResponse := dto.UserResponseFromModel(createdUser)

	response.JSON(ctx, status, map[string]interface{}{"user": modelResponse, "message": "user created successfully"})
}
