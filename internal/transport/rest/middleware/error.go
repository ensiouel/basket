package middleware

import (
	"errors"
	"github.com/ensiouel/apperror"
	"github.com/ensiouel/apperror/codes"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"net/http"
)

func Error(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		logger := logger.With(
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
		)

		err := c.Errors[0]
		code := http.StatusInternalServerError

		var apperr *apperror.Error
		if errors.As(err.Err, &apperr) {
			switch apperr.Code {
			case codes.Internal:
				logger.Error("internal error",
					slog.String("error", apperr.Error()),
				)

				code = http.StatusInternalServerError
			case codes.BadRequest:
				code = http.StatusBadRequest
			case codes.NotFound:
				code = http.StatusNotFound
			default:
				logger.Warn("unexpected error",
					slog.String("error", apperr.Error()),
				)

				code = http.StatusTeapot
			}

			c.JSON(code, gin.H{"error": apperr})
		} else {
			logger.Warn("unexpected error",
				slog.String("error", err.Error()),
			)

			c.JSON(code, gin.H{"error": "internal server error"})
		}
	}
}
