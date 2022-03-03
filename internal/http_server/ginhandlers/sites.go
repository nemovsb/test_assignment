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
)

type SiteHandler struct {
	TTL     uint
	Timeout uint
	DB      storage.DB
	Logger  *zap.Logger
	CommonHandler
}

func NewSiteHandler(config *di.ConfigApp, db storage.DB, logger *zap.Logger) *SiteHandler {
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

	site, rows := h.DB.GetSiteByName(searchUrl)

	if rows != 0 && time.Since(site.UpdatedAt) <= time.Second*time.Duration(h.TTL) {
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

	h.DB.CreateSite(searchUrl, duration)

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

	report, _ := h.DB.GetReportByDate(param.From, param.To)

	ctx.JSON(http.StatusOK, report)

}
