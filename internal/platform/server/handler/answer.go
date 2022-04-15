package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/patriciabonaldy/bequest_challenge/internal/business"
	"github.com/pkg/errors"
	"net/http"
)

type request struct {
	ID    string `json:"id" binding:"required"`
	Event string
	Data  []byte
}

// AnswerHandler returns an HTTP handler for answer creation.
func AnswerHandler(service business.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req request
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}

		err := service.CreateAnswer(ctx, req.Data)

		if err != nil {
			switch {
			case errors.Is(err, mooc.ErrInvalidCourseID),
				errors.Is(err, mooc.ErrEmptyCourseName), errors.Is(err, mooc.ErrInvalidCourseID):
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			default:
				ctx.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}

		ctx.Status(http.StatusCreated)
	}
}
