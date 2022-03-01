package ginhandlers

import (
	"fmt"
	"net/http"
	"time"

	"test_assignment/internal/storage"

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
	url, ok := ctx.GetQuery("search")
	if !ok {
		h.StatusBadRequest(ctx, fmt.Errorf("search param not found:"))
		return
	}

	// site := storage.Sites{}
	// h.DB.

	start := time.Now()
	// Код для измерения
	client := http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	// Создание
	h.DB.Create(&storage.Sites{Name: "D42", LoadingTime: time.Second})

	// Чтение
	var site storage.Sites
	h.DB.First(&site, 1) // find product with integer primary key

	fmt.Printf("site: %+v\n")

	_, err := client.Get(url)
	if err != nil {
		h.StatusInternalServerError(ctx, err)
		return
	}

	duration := time.Since(start)
	fmt.Printf("duration: %s\n", duration)

}

func (h *SiteHandler) GetReport(ctx *gin.Context) {

}
