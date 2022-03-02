package ginhandlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"test_assignment/internal/configuration/di"
	"test_assignment/internal/storage"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SiteHandler struct {
	TTL     uint
	Timeout uint
	DB      *gorm.DB
	Logger  *zap.Logger
	CommonHandler
}

func NewSiteHandler(config di.ConfigApp, db *gorm.DB, logger *zap.Logger) *SiteHandler {
	return &SiteHandler{
		TTL:     config.HttpServer.TTL,
		Timeout: config.HttpServer.Timeout,
		DB:      db,
		Logger:  logger,
	}
}

func (h *SiteHandler) CheckSite(ctx *gin.Context) {
	searchUrl, ok := ctx.GetQuery("search")
	if !ok {
		h.StatusBadRequest(ctx, errors.New(`"search" param not found`))
		return
	}

	_, err := url.ParseRequestURI(searchUrl)
	if err != nil {
		fmt.Printf("Error : %s\n", err)
		searchUrl = fmt.Sprint(`http://`, searchUrl)
	}

	fmt.Println("  url  :  ", searchUrl)

	var site storage.Sites
	res := h.DB.Where(`updated_at = (?)`, h.DB.Table("sites").Select(`max(updated_at)`).
		Group(`name`).Where("name = ?", searchUrl)).Find(&site)

	//fmt.Printf("res: %+v\n", res)
	//fmt.Printf("error: %+v\n", errors.Is(res.Error, gorm.ErrRecordNotFound))
	//fmt.Printf("site: %+v\n", site)
	// fmt.Printf("time.Now():  %+v\n", time.Now().UTC())
	// fmt.Printf("site.UpdatedAt:  %+v\n", site.UpdatedAt)
	// fmt.Printf("time.Since(site.UpdatedAt): %+v\n", time.Since(site.UpdatedAt))

	if res.RowsAffected != 0 && time.Since(site.UpdatedAt) <= time.Second*time.Duration(h.TTL) {
		ctx.JSON(http.StatusOK, gin.H{"duration": site.LoadingTime})
		return
	}

	start := time.Now()

	client := http.Client{
		Timeout: time.Duration(h.Timeout) * time.Second,
	}
	_, err = client.Get(searchUrl)
	if err != nil {
		h.StatusInternalServerError(ctx, err)
		return
	}

	duration := time.Since(start)
	fmt.Printf("duration: %s\n", duration)

	h.DB.Create(&storage.Sites{Name: searchUrl, LoadingTime: duration})

	ctx.JSON(http.StatusOK, gin.H{"duration": duration})

}

func (h *SiteHandler) GetReport(ctx *gin.Context) {

	type Param struct {
		From time.Time `form:"from" binding:"required"`
		To   time.Time `form:"to" binding:"required"`
	}

	param := new(Param)
	if err := ctx.ShouldBind(param); err != nil {
		h.StatusBadRequest(ctx, err)
		return
	}

	report := []storage.Report{}

	h.DB.Table("sites").
		Select(`name, avg(loading_time)::bigint AS "avg_duration"`).
		Where("created_at >= ?", param.From).
		Where("created_at <= ?", param.To).
		Group("name").
		Find(&report)
	//fmt.Printf("Report : %+v\n", report)

	ctx.JSON(http.StatusOK, report)

}
