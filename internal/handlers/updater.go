package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
)

func (h *Handler) updateGauge(c *gin.Context) {

	metricName := c.Param("metric")
	metricValue := c.Param("value")

	reg, err := regexp.Compile(`[.\d]`)
	if err != nil {
		logrus.Errorf("error while regex compiling: %s", err.Error())
	}
	if !reg.MatchString(metricValue) {
		h.statusBadRequest(c)
	} else {
		h.repos.Store(metricName, metricValue, false)
		c.Status(http.StatusOK)
	}
}

func (h *Handler) updateCounter(c *gin.Context) {

	metricName := c.Param("metric")
	metricValue := c.Param("value")

	reg, err := regexp.Compile(`[.\d]`)
	if err != nil {
		logrus.Errorf("error while regex compiling: %s", err.Error())
	}
	if !reg.MatchString(metricValue) {
		h.statusBadRequest(c)
	} else {
		h.repos.Store(metricName, metricValue, true)
		c.Status(http.StatusOK)
	}
}

func (h *Handler) singleGauge(c *gin.Context) {
	metricName := c.Param("metric")

	metricValue, err := h.repos.SingleMetric(metricName, false)
	if err != nil {
		h.statusNotFound(c)
		//errorResponse(c, http.StatusNotFound, "not found")
		return
	}
	_, err = c.Writer.WriteString(metricValue)
	c.Status(http.StatusOK)
	if err != nil {
		logrus.Printf("error while writing value %s of metric %s to body", metricValue, metricName)
	}
}

func (h *Handler) singleCounter(c *gin.Context) {
	metricName := c.Param("metric")

	metricValue, err := h.repos.SingleMetric(metricName, true)
	if err != nil {
		h.statusNotFound(c)
		//errorResponse(c, http.StatusNotFound, "not found")
		return
	}
	_, err = c.Writer.WriteString(metricValue)
	c.Status(http.StatusOK)
	if err != nil {
		logrus.Printf("error while writing value %s of metric %s to body", metricValue, metricName)
	}
}

func (h *Handler) startPage(c *gin.Context) {
	currentMetrics := h.repos.Current()
	c.HTML(http.StatusOK, "start_page.html", gin.H{
		"content": `"Welcome!"
There are some metrics for you.`,
		"Metrics": currentMetrics,
	})
}
