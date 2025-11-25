package handler

import (
	"avito-shop-service/internal/domain/purchase"
	"avito-shop-service/internal/domain/transfer"
	"avito-shop-service/internal/domain/user"
	"log/slog"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	log             *slog.Logger
	userHandler     user.UserHandler
	purchaseHandler purchase.PurchaseHandler
	transferHandler transfer.TransferHandler
	engine          *gin.Engine
}

func NewHandler(
	log *slog.Logger,
	userHandler user.UserHandler,
	purchaseHandler purchase.PurchaseHandler,
	transferHandler transfer.TransferHandler,
) *Handler {
	h := &Handler{
		log:             log,
		userHandler:     userHandler,
		purchaseHandler: purchaseHandler,
		transferHandler: transferHandler,
		engine:          gin.New(),
	}

	h.initMiddlewares()
	h.initRoutes()

	return h
}

func (h *Handler) NewEngine() *gin.Engine {
	return h.engine
}

func (h *Handler) initMiddlewares() {
	h.engine.Use(gin.Recovery())
	h.engine.Use(gin.Logger())
}

func (r *Handler) initRoutes() {
	apiV1 := r.engine.Group("/api/v1")

	authGroup := apiV1.Group("/auth")
	{
		authGroup.POST("/signup", r.userHandler.SignUpRoute)
		authGroup.POST("/signin", r.userHandler.SignInRoute)
	}

	protectedGroup := apiV1.Group("")
	{
		// /api/v1/transfers
		transfersGroup := protectedGroup.Group("/transfers")
		{
			transfersGroup.POST("", r.transferHandler.SendCoinsRoute)
			transfersGroup.GET("", r.transferHandler.GetHistoryRoute)
		}
		// /api/v1/purchases
		purchasesGroup := protectedGroup.Group("/purchases")
		{
			purchasesGroup.POST("", r.purchaseHandler.BuyItemRoute)
			purchasesGroup.GET("", r.purchaseHandler.GetHistoryRoute)
		}
	}
}
