package infrahttp

import (
	"time"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Error     interface{} `json:"error,omitempty"`
	Timestamp string      `json:"timestamp"`
}

func SendSuccessResponse(gc *gin.Context, statusCode int, message string, data interface{}) {
	gc.JSON(statusCode, BaseResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

func SendErrorResponse(gc *gin.Context, statusCode int, message string, err interface{}) {
	gc.AbortWithStatusJSON(statusCode, BaseResponse{
		Success:   false,
		Message:   message,
		Error:     err,
		Timestamp: time.Now().Format(time.RFC3339),
	})
}
