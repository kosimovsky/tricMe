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

func errorsResponseWithStatusCode(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errors{Message: message})
}

func errorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.Status(statusCode)
}

func (h *Handler) statusNotImplemented(c *gin.Context) {
	errorResponse(c, http.StatusNotImplemented, "Not implemented")
}

func (h *Handler) statusNotImplementedRegex(c *gin.Context) {
	reg, err := regexp.Compile(`(gauge\b|counter\b)`)
	if err != nil {
		logrus.Errorf("error while regex compiling: %s", err.Error())
	}
	metricType := c.Param("regex")
	if !reg.MatchString(metricType) {
		h.statusNotImplemented(c)
	}
}

func (h *Handler) statusNotFound(c *gin.Context) {
	errorResponse(c, http.StatusNotFound, "Not found")
}

func (h *Handler) statusBadRequest(c *gin.Context) {
	errorResponse(c, http.StatusBadRequest, "Bad request")
}
