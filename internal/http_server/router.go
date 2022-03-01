package http_server

import (
	"github.com/gin-gonic/gin"
)

type SiteHandler interface {
	CheckSite(ctx *gin.Context)
	GetReport(ctx *gin.Context)
}

type HandlerSet struct {
	SiteHandler
}

func NewHandlerSet(site SiteHandler) HandlerSet {
	return HandlerSet{
		SiteHandler: site,
	}
}

func NewRouter(h HandlerSet) (router *gin.Engine) {
	router = gin.Default()

	sites := router.Group("/sites")
	{
		sites.GET("/", h.SiteHandler.CheckSite)
		sites.GET("/stats", h.SiteHandler.GetReport)
	}

	return router
}
