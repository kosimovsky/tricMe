package handlers

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"

	"github.com/kosimovsky/tricMe/internal/storage"
)

type Handler struct {
	keeper storage.Storekeeper
}

func NewHandler(keeper storage.Storekeeper) *Handler {
	return &Handler{keeper: keeper}
}

func (h *Handler) MetricsRouter() *gin.Engine {
	router := gin.New()

	htmlStart := router.Group("/")
	htmlStart.Use(gzip.Gzip(gzip.DefaultCompression))
	{
		router.LoadHTMLGlob("templates/*.html")
		htmlStart.GET("/", h.startPage)
	}

	update := router.Group("/update")
	{
		update.POST("", h.statusNotImplemented)
		update.Use(gzip.Gzip(gzip.DefaultCompression))
		update.POST("/", h.updateMetric)
		r := update.Group("/:regex", h.statusNotImplementedRegex)
		{
			r.POST("", h.statusNotImplementedRegex)
			mType := r.Group("/:metric", h.statusNotImplemented)
			{
				mType.POST("", h.statusNotImplemented)
				mType.POST("/:value", h.statusNotImplemented)
			}
		}
		gauge := update.Group("/gauge")
		{
			gauge.POST("", h.statusNotFound)
			gauge.POST("/", h.statusNotFound)
			metric := gauge.Group("/:metric", h.statusNotFound)
			{
				metric.POST("/:value", h.updateGauge)
			}
		}
		counter := update.Group("/counter")
		{
			counter.POST("", h.statusNotFound)
			counter.POST("/", h.statusNotFound)
			metric := counter.Group("/:metric", h.statusNotFound)
			{
				metric.POST("/:value", h.updateCounter)
			}
		}
	}

	value := router.Group("/value")
	{
		value.GET("", h.statusNotImplemented)
		value.POST("/", h.valueOf)
		value.Use(gzip.Gzip(gzip.DefaultCompression))
		r := value.Group("/:regex", h.statusNotImplementedRegex)
		{
			r.GET("", h.statusNotImplementedRegex)
			mType := r.Group("/:metric", h.statusNotImplemented)
			{
				mType.GET("", h.statusNotImplemented)
			}
		}
		gauge := value.Group("/gauge")
		{
			gauge.GET("", h.statusNotFound)
			gauge.GET("/:metric", h.singleGauge)
		}
		counter := value.Group("/counter")
		{
			counter.GET("", h.statusNotFound)
			counter.GET("/:metric", h.singleCounter)
		}
	}
	return router
}
