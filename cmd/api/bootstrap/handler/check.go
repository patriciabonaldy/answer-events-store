package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// check godoc
// @Tags         CheckHandler
// @Accept       json
// @Produce      json
// @Success      200  {string}  string  "pong"
// @Failure      400  {string}  string  "ok"
// @Failure      500  {string}  string  "ok"
// @Router       /health [get]
func CheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	}
}
