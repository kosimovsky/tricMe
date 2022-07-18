package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) updateGauge(c *gin.Context) {

	metricName := c.Param("metric")
	metricValue := c.Param("value")

	err := h.repos.Store(metricName, metricValue, false)
	c.Status(http.StatusOK)
	if err != nil {
		logrus.Printf("error while storing metric %s with value %s", metricName, metricValue)
	}
}

func (h *Handler) updateCounter(c *gin.Context) {

	metricName := c.Param("metric")
	metricValue := c.Param("value")

	err := h.repos.Store(metricName, metricValue, true)
	c.Status(http.StatusOK)
	if err != nil {
		logrus.Printf("error while storing metric %s with value %s", metricName, metricValue)
	}
}

func (h *Handler) singleGauge(c *gin.Context) {
	metricName := c.Param("metric")

	metricValue, err := h.repos.SingleMetric(metricName, false)
	if err != nil {
		errorResponse(c, http.StatusNotFound, "not found")
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
		errorResponse(c, http.StatusNotFound, "not found")
		return
	}
	_, err = c.Writer.WriteString(metricValue)
	c.Status(http.StatusOK)
	if err != nil {
		logrus.Printf("error while writing value %s of metric %s to body", metricValue, metricName)
	}
}

func (h *Handler) startPage(c *gin.Context) {
	c.HTML(http.StatusOK, "start_page.html", gin.H{
		"content": "this is the start page",
	})

	//m, _ := h.repos.Marshal()

}
