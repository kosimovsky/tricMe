package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"regexp"
)

func (h *Handler) validateValue(c *gin.Context) {
	value := c.Param("value")

	reg, err := regexp.Compile(`[.\d]`)
	if err != nil {
		logrus.Errorf("error while regex compiling: %s", err.Error())
	}
	if !reg.MatchString(value) {
		h.statusBadRequest(c)
		return
	}
	h.updateCounter(c)
}
