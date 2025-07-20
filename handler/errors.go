package handler

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func JSONError(ctx *gin.Context, code int, msg string) {
	zap.L().Error("HTTP error", zap.Int("status", code), zap.String("message", msg))
	ctx.IndentedJSON(code, gin.H{"message": msg})
}
