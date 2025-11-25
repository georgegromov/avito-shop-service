package handler

import (
	"avito-shop-service/internal/domain/transfer"
	"avito-shop-service/internal/domain/wallet"
	infrahttp "avito-shop-service/internal/infra/http"
	"avito-shop-service/pkg/utils"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transferHandler struct {
	log             *slog.Logger
	validator       infrahttp.Validator
	transferService transfer.TransferService
}

func NewTransferHandler(
	log *slog.Logger,
	validator infrahttp.Validator,
	transferService transfer.TransferService,
) transfer.TransferHandler {
	return &transferHandler{
		log:             log,
		validator:       validator,
		transferService: transferService,
	}
}

// SendCoinsRoute implements transfer.TransferHandler.
func (h *transferHandler) SendCoinsRoute(gc *gin.Context) {
	const op = "handler.transferHandler.SendCoinsRoute"

	userCtx, err := utils.GetUserFromCtx(gc)
	if err != nil {
		infrahttp.SendErrorResponse(gc, http.StatusUnauthorized, infrahttp.ErrUnauthorized.Error(), gin.H{"error": infrahttp.ErrUnauthorized.Error()})
		return
	}

	ctx := gc.Request.Context()
	log := h.log.With(slog.String("op", op))

	var reqBody transfer.SendCoinsRequestDTO
	if err := gc.ShouldBindJSON(&reqBody); err != nil {
		log.Error("an error occurred while binding request body", slog.Any("err", err))
		infrahttp.SendErrorResponse(gc, http.StatusBadRequest, infrahttp.ErrBadRequest.Error(), gin.H{"error": infrahttp.ErrBadRequest.Error()})
		return
	}

	if err := h.validator.Validate(ctx, reqBody); err != nil {
		log.Error("an error occurred while validating request body", slog.Any("err", err))
		infrahttp.SendErrorResponse(gc, http.StatusBadRequest, infrahttp.ErrBadRequest.Error(), gin.H{"error": infrahttp.ErrBadRequest.Error()})
		return
	}

	if err := h.transferService.SendCoins(ctx, userCtx.ID, reqBody.ToUserID, reqBody.Amount); err != nil {
		log.Error("failed to send coins", slog.Any("err", err))
		if errors.Is(err, wallet.ErrInsufficientFunds) {
			infrahttp.SendErrorResponse(gc, http.StatusPaymentRequired, wallet.ErrInsufficientFunds.Error(), gin.H{"error": wallet.ErrInsufficientFunds.Error()})
			return
		}
		if errors.Is(err, transfer.ErrCannotSendToYourself) {
			infrahttp.SendErrorResponse(gc, http.StatusBadRequest, transfer.ErrCannotSendToYourself.Error(), gin.H{"error": transfer.ErrCannotSendToYourself.Error()})
			return
		}
		infrahttp.SendErrorResponse(gc, http.StatusInternalServerError, infrahttp.ErrInternalServer.Error(), gin.H{"error": infrahttp.ErrInternalServer.Error()})
		return
	}

	infrahttp.SendSuccessResponse(gc, http.StatusOK, "coins sent successfully", nil)
}

// GetHistoryRoute implements transfer.TransferHandler.
func (t *transferHandler) GetHistoryRoute(gc *gin.Context) {
	panic("unimplemented")
}
