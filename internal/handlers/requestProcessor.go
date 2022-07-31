package handlers

import (
	"encoding/json"
	tricme "github.com/kosimovsky/tricMe"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) updateMetric(c *gin.Context) {
	m := new(tricme.Metrics)

	if err := json.NewDecoder(c.Request.Body).Decode(m); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
	}
	h.keeper.Store(*m)
	c.Status(http.StatusOK)
}

func (h *Handler) valueOf(c *gin.Context) {
	bodyData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.Status(http.StatusBadRequest)
	}
	var inMetric tricme.Metrics

	err = json.Unmarshal(bodyData, &inMetric)
	if err != nil {
		logrus.Errorf("error while unmarsalling body to struct: %s", err.Error())
	}
	outMetric, err := h.keeper.SingleMetric(inMetric.ID, inMetric.MType)
	if err != nil {
		logrus.Println(err.Error())
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, outMetric)
}

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
		metric := convert("gauge", metricName, metricValue)
		h.keeper.Store(metric)
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
		metric := convert("counter", metricName, metricValue)
		h.keeper.Store(metric)
		c.Status(http.StatusOK)
	}
}

func (h *Handler) singleGauge(c *gin.Context) {
	metricName := c.Param("metric")
	mType := strings.Split(c.FullPath(), "/")[2]

	metric, err := h.keeper.SingleMetric(metricName, mType)
	if err != nil {
		h.statusNotFound(c)
		//errorResponse(c, http.StatusNotFound, "not found")
		return
	}
	_, err = c.Writer.WriteString(strconv.FormatFloat(*metric.Value, 'f', -1, 64))
	c.Status(http.StatusOK)
	if err != nil {
		logrus.Printf("error while writing value of metric %s to body", metricName)
	}
}

func (h *Handler) singleCounter(c *gin.Context) {
	metricName := c.Param("metric")
	mType := strings.Split(c.FullPath(), "/")[2]

	metric, err := h.keeper.SingleMetric(metricName, mType)
	if err != nil {
		h.statusNotFound(c)
		//errorResponse(c, http.StatusNotFound, "not found")
		return
	}
	_, err = c.Writer.WriteString(strconv.FormatInt(*metric.Delta, 10))
	c.Status(http.StatusOK)
	if err != nil {
		logrus.Printf("error while writing value of metric %s to body", metricName)
	}
}

func (h *Handler) startPage(c *gin.Context) {
	currentMetrics := h.keeper.Current()
	c.HTML(http.StatusOK, "start_page.html", gin.H{
		"content": `Welcome!
There are some metrics for you.`,
		"Metrics": currentMetrics,
	})
}

func convert(mType, id, content string) tricme.Metrics {
	m := new(tricme.Metrics)

	m.ID = id
	m.MType = mType
	if mType == "counter" {
		delta, _ := strconv.ParseInt(content, 10, 64)
		m.Delta = &delta
		m.Value = nil
	} else {
		value, _ := strconv.ParseFloat(content, 64)
		m.Value = &value
		m.Delta = nil
	}
	return *m
}
