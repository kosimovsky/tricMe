package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
)

type errors struct {
	Message string `json:"message"`
}

func errorsResponseWithStatusCode(ctx *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	ctx.AbortWithStatusJSON(statusCode, errors{Message: message})
}

func errorResponse(ctx *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	ctx.Status(statusCode)
}

func (h *Handler) statusNotImplemented(ctx *gin.Context) {
	errorResponse(ctx, http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) statusNotImplementedRegex(ctx *gin.Context) {
	reg, err := regexp.Compile(`(gauge\b|counter\b)`)
	if err != nil {
		logrus.Errorf("error while regex compiling: %s", err.Error())
	}
	metricType := ctx.Param("regex")
	if !reg.MatchString(metricType) {
		h.statusNotImplemented(ctx)
	}
}

func (h *Handler) statusNotFound(ctx *gin.Context) {
	errorResponse(ctx, http.StatusNotFound, "Not found")
}
