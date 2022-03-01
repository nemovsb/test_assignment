package ginhandlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
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
	serchUrl, ok := ctx.GetQuery("search")
	if !ok {
		h.StatusBadRequest(ctx, errors.New("search param not found"))
		return
	}

	_, err := url.ParseRequestURI(serchUrl)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		serchUrl = fmt.Sprint(`http://`, serchUrl)
	}

	fmt.Println("  url  :  ", serchUrl)

	var site storage.Sites
	res := h.DB.Find(&site, storage.Sites{Name: serchUrl})
	//fmt.Printf("res: %+v\n", res)
	//fmt.Printf("error: %+v\n", errors.Is(res.Error, gorm.ErrRecordNotFound))
	//fmt.Printf("site: %+v\n", site)
	fmt.Printf("time.Now():  %+v\n", time.Now().UTC())
	fmt.Printf("site.UpdatedAt:  %+v\n", site.UpdatedAt)
	fmt.Printf("time.Since(site.UpdatedAt): %+v\n", time.Since(site.UpdatedAt))

	if res.RowsAffected != 0 && time.Since(site.UpdatedAt) <= time.Second*time.Duration(30) {
		ctx.JSON(http.StatusOK, gin.H{"duration": site.LoadingTime})
		return
	}

	start := time.Now()
	//Код для измерения
	client := http.Client{
		Timeout: time.Duration(60) * time.Second,
	}

	_, err = client.Get(serchUrl)
	if err != nil {
		h.StatusInternalServerError(ctx, err)
		return
	}

	duration := time.Since(start)
	fmt.Printf("duration: %s\n", duration)

	h.DB.Create(&storage.Sites{Name: serchUrl, LoadingTime: duration})

	ctx.JSON(http.StatusOK, gin.H{"duration": duration})

}

func (h *SiteHandler) GetReport(ctx *gin.Context) {

}
