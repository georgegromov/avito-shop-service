package handler

import (
	"avito-shop-service/internal/domain/item"
	"avito-shop-service/internal/domain/purchase"
	"avito-shop-service/internal/domain/wallet"
	infrahttp "avito-shop-service/internal/infra/http"
	"avito-shop-service/pkg/utils"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type purchaseHandler struct {
	log             *slog.Logger
	validator       infrahttp.Validator
	purchaseService purchase.PurchaseService
}

func NewPurchaseHandler(
	log *slog.Logger,
	validator infrahttp.Validator,
	purchaseService purchase.PurchaseService,
) purchase.PurchaseHandler {
	return &purchaseHandler{
		log:             log,
		validator:       validator,
		purchaseService: purchaseService,
	}
}

// BuyItemRoute implements purchase.PurchaseHandler.
func (h *purchaseHandler) BuyItemRoute(gc *gin.Context) {
	const op = "handler.purchaseHandler.BuyItemRoute"

	userCtx, err := utils.GetUserFromCtx(gc)
	if err != nil {
		infrahttp.SendErrorResponse(gc, http.StatusUnauthorized, infrahttp.ErrUnauthorized.Error(), gin.H{"error": infrahttp.ErrUnauthorized.Error()})
		return
	}

	ctx := gc.Request.Context()
	log := h.log.With(slog.String("op", op))

	var reqBody item.BuyItemRequestDTO
	if err := gc.ShouldBindJSON(&reqBody); err != nil {
		log.Error("an error occurred while binding request body", slog.Any("err", err))
		infrahttp.SendErrorResponse(gc, http.StatusBadRequest, infrahttp.ErrBadRequest.Error(), gin.H{"error": infrahttp.ErrBadRequest.Error()})
		return
	}

	if err := h.validator.Validate(ctx, reqBody); err != nil {
		log.Error("an error occurred while validating request body", slog.Any("err", err))
		infrahttp.SendErrorResponse(
			gc,
			http.StatusBadRequest,
			infrahttp.ErrBadRequest.Error(),
			gin.H{"error": infrahttp.ErrBadRequest.Error()},
		)
		return
	}

	if err := h.purchaseService.BuyItem(ctx, userCtx.ID, reqBody.ItemID, reqBody.Quantity); err != nil {
		log.Error("an error occurred while buying item", slog.Any("err", err))
		if errors.Is(err, wallet.ErrInsufficientFunds) {
			infrahttp.SendErrorResponse(gc, http.StatusPaymentRequired, wallet.ErrInsufficientFunds.Error(), gin.H{"error": wallet.ErrInsufficientFunds.Error()})
			return
		}
		infrahttp.SendErrorResponse(gc, http.StatusInternalServerError, infrahttp.ErrInternalServer.Error(), gin.H{"error": infrahttp.ErrInternalServer.Error()})
		return
	}

	infrahttp.SendSuccessResponse(gc, http.StatusOK, "purchase completed successfully", nil)
}

// GetHistoryRoute implements purchase.PurchaseHandler.
func (h *purchaseHandler) GetHistoryRoute(gc *gin.Context) {
	const op = "handler.purchaseHandler.GetHistoryRoute"

	userCtx, err := utils.GetUserFromCtx(gc)
	if err != nil {
		fmt.Printf("no user set in context: %v", err)
		infrahttp.SendErrorResponse(gc, http.StatusUnauthorized, infrahttp.ErrUnauthorized.Error(), gin.H{"error": infrahttp.ErrUnauthorized.Error()})
		return
	}

	ctx := gc.Request.Context()
	log := h.log.With(slog.String("op", op))

	history, err := h.purchaseService.GetHistory(ctx, userCtx.ID)
	if err != nil {
		log.Error("failed to get purchase history", slog.Any("err", err))
		infrahttp.SendErrorResponse(gc, http.StatusInternalServerError, infrahttp.ErrInternalServer.Error(), gin.H{"error": infrahttp.ErrInternalServer.Error()})
		return
	}

	infrahttp.SendSuccessResponse(gc, http.StatusOK, "purchase history retrieved successfully", history)
}
