package ginhandlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SiteHandler struct {
	TTL     uint
	Timeout uint
	Metrics
	Cache  Cacher
	DB     DB
	Logger *zap.Logger
	CommonHandler
}

type Cacher interface {
	Get(string) (time.Duration, bool)
	Set(searchUrlname string, duration time.Duration)
}

type DB interface {
	GetSiteDuration(name string) (time.Duration, bool)
	CreateSite(name string, duration time.Duration) int64
	GetReportByDate(from, to time.Time) (*[]Report, int64)
}

type Report struct {
	Name        string        `json:"name"`
	AvgDuration time.Duration `json:"avg_duration"`
}

func NewSiteHandler(ttl uint, timeout uint, cache Cacher, db DB, logger *zap.Logger) *SiteHandler {
	return &SiteHandler{
		TTL:     ttl,
		Timeout: timeout,
		Metrics: *NewMetrics(),
		Cache:   cache,
		DB:      db,
		Logger:  logger,
	}
}

func (h *SiteHandler) CheckSite(ctx *gin.Context) {

	h.Metrics.RequestCount.Inc()

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

	siteDuration, ok := h.Cache.Get(searchUrl)
	if ok {
		ctx.JSON(http.StatusOK, gin.H{"duration": siteDuration})
		return
	}

	siteDuration, ok = h.DB.GetSiteDuration(searchUrl)

	if ok {
		ctx.JSON(http.StatusOK, gin.H{"duration": siteDuration})
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func(url string, dur time.Duration) {
		h.Cache.Set(url, dur)
	}(searchUrl, duration)

	h.DB.CreateSite(searchUrl, duration)

	ctx.JSON(http.StatusOK, gin.H{"duration": duration})

}

func (h *SiteHandler) GetReport(ctx *gin.Context) {

	h.Metrics.RequestCount.Inc()

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
