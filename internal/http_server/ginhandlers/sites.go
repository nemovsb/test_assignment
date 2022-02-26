package ginhandlers

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SiteHandler struct {
	DB     *gorm.DB
	Logger *zap.Logger
	CommonHandler
}

func NewSiteHandler(db *gorm.DB, logger *zap.Logger) *SiteHandler {
	return &SiteHandler{
		DB:     db,
		Logger: logger,
	}
}

func (h *SiteHandler) CheckSite(ctx *gin.Context) {

}

func (h *SiteHandler) GetReport(ctx *gin.Context) {

}
