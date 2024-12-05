package routes

import (
	"net/http"

	"github.com/bakare-dev/simple-bank-api/internal/user/dto"
	"github.com/bakare-dev/simple-bank-api/internal/user/model"
	"github.com/bakare-dev/simple-bank-api/internal/user/service"
	"github.com/bakare-dev/simple-bank-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) HandleCreateUser(ctx *gin.Context) {
	var createUserDTO dto.CreateUserRequest
	if err := ctx.ShouldBindJSON(&createUserDTO); err != nil {
		response.Error(ctx, http.StatusUnprocessableEntity, nil, "Invalid request data")
		return
	}

	user := createUserDTO.ToModel()
	createdUser, status, err := h.userService.CreateUser(ctx.Request.Context(), user)
	if err != nil {
		response.Error(ctx, status, err, err.Error())
		return
	}

	modelResponse := dto.UserResponseFromModel(createdUser)
	response.JSON(ctx, status, modelResponse, "User created successfully, check mail for otp")
}

func (h *UserHandler) HandleCreateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, nil, "Unauthorized: User ID not found")
		return
	}

	role, exists := ctx.Get("role")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, nil, "Unauthorized: Role not found")
		return
	}

	if role != "admin" && role != "customer" {
		response.Error(ctx, http.StatusForbidden, nil, "Insufficient permissions")
		return
	}

	var createProfileDto dto.CreateProfileRequest
	if err := ctx.ShouldBindJSON(&createProfileDto); err != nil {
		response.Error(ctx, http.StatusUnprocessableEntity, nil, "Invalid request data")
		return
	}

	profile := createProfileDto.ToModel()
	profile.UserID = userID.(string)

	_, status, err := h.userService.CreateProfile(ctx.Request.Context(), profile)
	if err != nil {
		response.Error(ctx, status, nil, err.Error())
		return
	}

	response.JSON(ctx, status, nil, "Profile created successfully")
}

func (h *UserHandler) HandleLogin(ctx *gin.Context) {
	var loginDto dto.LoginRequest
	if err := ctx.ShouldBindJSON(&loginDto); err != nil {
		response.Error(ctx, http.StatusUnprocessableEntity, nil, "Invalid request data")
		return
	}

	user := loginDto.ToModel()
	authenticatedUser, status, token, profile, err := h.userService.SignIn(ctx.Request.Context(), user)
	if err != nil {
		if authenticatedUser != nil {
			response.JSON(ctx, status, gin.H{"userId": authenticatedUser}, err.Error())
			return
		}
		response.Error(ctx, status, nil, err.Error())
		return
	}

	response.JSON(ctx, status, gin.H{
		"user":    authenticatedUser,
		"token":   token,
		"profile": profile,
	}, "Login successful")
}

func (h *UserHandler) HandleGetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, nil, "Unauthorized: User ID not found")
		return
	}

	role, exists := ctx.Get("role")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, nil, "Unauthorized: Role not found")
		return
	}

	if role != "admin" && role != "customer" {
		response.Error(ctx, http.StatusForbidden, nil, "Insufficient permissions")
		return
	}

	user := &model.User{
		ID: userID.(string),
	}

	profileData, status, err := h.userService.GetProfile(ctx.Request.Context(), user)
	if err != nil {
		response.Error(ctx, status, err, err.Error())
		return
	}

	response.JSON(ctx, status, profileData, "Profile retrieved successfully")
}

func (h *UserHandler) HandleActivateUser(ctx *gin.Context) {
	var activateUserDto dto.ActivateUserRequest

	if err := ctx.ShouldBindJSON(&activateUserDto); err != nil {
		response.Error(ctx, http.StatusUnprocessableEntity, nil, "Invalid request data")
		return
	}

	statusCode, err := h.userService.ActivateUser(ctx, &model.User{ID: activateUserDto.ID}, activateUserDto.Otp)

	if err != nil {
		response.Error(ctx, statusCode, nil, err.Error())
		return
	}

	response.JSON(ctx, http.StatusOK, nil, "User activated successfully")
}

func (h *UserHandler) HandleUpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, nil, "Unauthorized: User ID not found")
		return
	}

	role, exists := ctx.Get("role")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, nil, "Unauthorized: Role not found")
		return
	}

	if role != "admin" && role != "customer" {
		response.Error(ctx, http.StatusForbidden, nil, "Insufficient permissions")
		return
	}

	var updateProfileDTO dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&updateProfileDTO); err != nil {
		response.Error(ctx, http.StatusUnprocessableEntity, err, "Invalid request data")
		return
	}

	profile := model.Profile{
		UserID: userID.(string),
	}

	if updateProfileDTO.FirstName != "" {
		profile.FirstName = updateProfileDTO.FirstName
	}

	if updateProfileDTO.LastName != "" {
		profile.LastName = updateProfileDTO.LastName
	}

	if updateProfileDTO.PhoneNumber != "" {
		profile.PhoneNumber = updateProfileDTO.PhoneNumber
	}

	if profile.FirstName == "" && profile.LastName == "" && profile.PhoneNumber == "" {
		response.Error(ctx, http.StatusUnprocessableEntity, nil, "No fields provided to update")
		return
	}

	res, status, err := h.userService.UpdateProfile(ctx.Request.Context(), &profile)
	if err != nil {
		response.Error(ctx, status, err, err.Error())
		return
	}

	response.JSON(ctx, status, res, "Profile updated successfully")
}

func (h *UserHandler) HandleResendOtp(ctx *gin.Context) {
	var resendOtpDto dto.ResendOtpRequest
	if err := ctx.ShouldBindJSON(&resendOtpDto); err != nil {
		response.Error(ctx, http.StatusUnprocessableEntity, err, "Invalid request data")
		return
	}

	res, status, err := h.userService.ResendOtp(ctx.Request.Context(), resendOtpDto.ID)
	if err != nil {
		response.Error(ctx, status, err, err.Error())
		return
	}

	response.JSON(ctx, status, res, "otp sent")
}

func (h *UserHandler) HandleDeactivateUser(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, nil, "Unauthorized: User ID not found")
		return
	}

	role, exists := ctx.Get("role")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, nil, "Unauthorized: Role not found")
		return
	}

	if role != "admin" && role != "customer" {
		response.Error(ctx, http.StatusForbidden, nil, "Insufficient permissions")
		return
	}

	statusCode, err := h.userService.DeactivateUser(ctx, userID.(string))

	if err != nil {
		response.Error(ctx, statusCode, nil, err.Error())
		return
	}

	response.JSON(ctx, http.StatusOK, nil, "User deactivated successfully")
}

func (h *UserHandler) HandleChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, nil, "Unauthorized: User ID not found")
		return
	}

	role, exists := ctx.Get("role")
	if !exists {
		response.Error(ctx, http.StatusUnauthorized, nil, "Unauthorized: Role not found")
		return
	}

	if role != "admin" && role != "customer" {
		response.Error(ctx, http.StatusForbidden, nil, "Insufficient permissions")
		return
	}

	var changePasswordDto dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&changePasswordDto); err != nil {
		response.Error(ctx, http.StatusUnprocessableEntity, err, "Invalid request data")
		return
	}

	statusCode, err := h.userService.ChangePassword(ctx, changePasswordDto.Password, userID.(string))

	if err != nil {
		response.Error(ctx, statusCode, nil, err.Error())
		return
	}

	response.JSON(ctx, http.StatusOK, nil, "password change")
}
