package api

import (
	"database/sql"
	"net/http"

	db "github.com/bakare-dev/simple-bank-api/db/sqlc"
	"github.com/bakare-dev/simple-bank-api/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type createUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email" usage:"Email should be in lowercase"`
	Password string `json:"password" binding:"required,min=8" usage:"Password must be at least 8 characters long"`
	Phoneno  string `json:"phoneno" binding:"required" usage:"Please provide a phone number"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)

	println(hashedPassword)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		Name:        req.Name,
		Email:       req.Email,
		Password:    hashedPassword,
		PhoneNumber: sql.NullString{String: req.Phoneno, Valid: req.Phoneno != ""},
	}

	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"userId": user.ID, "message": "User created successfully"})
}

type getUserRequest struct {
	Id string `form:"id" binding:"required"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req getUserRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request parameter"})
		return
	}

	parsedUUID, err := uuid.Parse(req.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	user, err := server.store.GetUser(ctx, parsedUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user.Password = ""

	ctx.JSON(http.StatusOK, user)
}

func (server *Server) updateUserpassword(ctx *gin.Context) {

}

type getUsersRequest struct {
	Page int32 `form:"page" binding:"required,min=1"`
	Size int32 `form:"size" binding:"required,min=5"`
}

func (server *Server) getUsers(ctx *gin.Context) {
	var req getUsersRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.Size,
		Offset: (req.Page - 1) * req.Size,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	for i := range users {
		users[i].Password = ""
	}

	ctx.JSON(http.StatusOK, users)
}

type deleteUserRequest struct {
	Id string `form:"id" binding:"required"`
}

func (server *Server) deleteUsers(ctx *gin.Context) {
	var req deleteUserRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	parsedUUID, err := uuid.Parse(req.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	err = server.store.DeleteUser(ctx, parsedUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error occurred deleting user", "error": err.Error()})
		return
	}

	_, err = server.store.GetUser(ctx, parsedUUID)
	if err == nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "user was not deleted"})
		return
	}

	if err != sql.ErrNoRows {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
