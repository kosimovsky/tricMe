package handlers

import (
	"compress/gzip"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
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
	} else {
		h.updateCounter(c)
	}
}

func (h *Handler) compressHandler(c *gin.Context) {
	if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		c.Next()
	} else {
		fmt.Println("here")

		gz, err := gzip.NewWriterLevel(c.Writer, gzip.DefaultCompression)
		if err != nil {
			logrus.Error(err.Error())
			return
		}
		defer gz.Close()
		c.Writer.Header().Set("Content-Encoding", "gzip")
		//c.Next()
	}
}

func (h *Handler) compressValueHandler(c *gin.Context) {
	if !strings.Contains(c.Request.Header.Get("Accept-Encoding"), "gzip") {
		c.Next()
	} else {
		fmt.Println("here")

		gz, err := gzip.NewWriterLevel(c.Writer, gzip.DefaultCompression)
		if err != nil {
			logrus.Error(err.Error())
			return
		}
		defer gz.Close()
		c.Writer.Header().Set("Content-Encoding", "gzip")
		//c.Next()
	}
}
