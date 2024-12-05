package routes

import (
	"net/http"

	"github.com/bakare-dev/simple-bank-api/internal/core/dto"
	"github.com/bakare-dev/simple-bank-api/internal/core/service"
	"github.com/bakare-dev/simple-bank-api/pkg/response"
	"github.com/gin-gonic/gin"
)

type CoreHandler struct {
	coreService *service.CoreService
}

func NewCoreHandler(coreService *service.CoreService) *CoreHandler {
	return &CoreHandler{coreService: coreService}
}

func (h *CoreHandler) HandleCreateAccount(ctx *gin.Context) {
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

	var createAccountDto dto.CreateAccountRequest
	if err := ctx.ShouldBindJSON(&createAccountDto); err != nil {
		response.Error(ctx, http.StatusUnprocessableEntity, err, "Invalid request data")
		return
	}

	account, err := createAccountDto.ToModel()

	if err != nil {
		response.Error(ctx, http.StatusUnprocessableEntity, err, "Invalid request data")
		return
	}

	statusCode, errMsg := h.coreService.CreateAccount(ctx, userID.(string), account)

	if errMsg != nil {
		response.Error(ctx, statusCode, nil, *errMsg)
		return
	}

	response.JSON(ctx, statusCode, nil, "Account created")
}

func (h *CoreHandler) HandleGetAccount(ctx *gin.Context) {
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

	data, statusCode, errMsg := h.coreService.GetAccount(ctx, userID.(string))

	if errMsg != nil {
		response.Error(ctx, statusCode, nil, *errMsg)
		return
	}

	response.JSON(ctx, statusCode, data, "Account retrieved successfully")
}

func (h *CoreHandler) HandleGetAccountByAccountNumber(ctx *gin.Context) {
	_, exists := ctx.Get("userID")
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

	var getAccountDto dto.GetAccountByAccountNumberRequest
	if err := ctx.ShouldBindQuery(&getAccountDto); err != nil {
		response.Error(ctx, http.StatusUnprocessableEntity, err, "Invalid request data")
		return
	}

	data, statusCode, errMsg := h.coreService.GetAccountByAccountNumber(ctx, getAccountDto.Number)

	if errMsg != nil {
		response.Error(ctx, statusCode, nil, *errMsg)
		return
	}

	response.JSON(ctx, statusCode, data, "Account retrieved successfully")
}

func (h *CoreHandler) HandleGetAccountBalance(ctx *gin.Context) {
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

	data, statusCode, errMsg := h.coreService.GetAccountBalance(ctx, userID.(string))

	if errMsg != nil {
		response.Error(ctx, statusCode, nil, *errMsg)
		return
	}

	response.JSON(ctx, statusCode, data, "Account balance retrieved")
}

func (h *CoreHandler) HandleGetAccountTransactions(ctx *gin.Context) {
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

	var accountTransactionDto dto.GetAccountTransactionRequest
	if err := ctx.ShouldBindQuery(&accountTransactionDto); err != nil {
		response.Error(ctx, http.StatusUnprocessableEntity, err, "Invalid request data")
		return
	}

	data, statusCode, errMsg := h.coreService.GetTransactions(ctx, userID.(string), accountTransactionDto.Page, accountTransactionDto.Size)

	if errMsg != nil {
		response.Error(ctx, statusCode, nil, *errMsg)
		return
	}

	response.JSON(ctx, statusCode, data, "Transactions retrieved")
}

func (h *CoreHandler) HandleGetTransaction(ctx *gin.Context) {
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

	var getTransactionDto dto.GetTransactionRequest
	if err := ctx.ShouldBindQuery(&getTransactionDto); err != nil {
		response.Error(ctx, http.StatusUnprocessableEntity, err, "Invalid request data")
		return
	}

	data, statusCode, errMsg := h.coreService.GetTransaction(ctx, userID.(string), getTransactionDto.TransactionId)

	if errMsg != nil {
		response.Error(ctx, statusCode, nil, *errMsg)
		return
	}

	response.JSON(ctx, statusCode, data, "Transaction retrieved")
}

func (h *CoreHandler) HandleTransfer(ctx *gin.Context) {
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

	var transferDto dto.TransferRequest
	if err := ctx.ShouldBindQuery(&transferDto); err != nil {
		response.Error(ctx, http.StatusUnprocessableEntity, err, "Invalid request data")
		return
	}

	statusCode, errMsg := h.coreService.Transfer(ctx, userID.(string), transferDto.ToAccountId, transferDto.Amount, transferDto.Pin)

	if errMsg != nil {
		response.Error(ctx, statusCode, nil, *errMsg)
		return
	}

	response.JSON(ctx, statusCode, nil, "Transfer successful")
}
