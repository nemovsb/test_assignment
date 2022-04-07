package ginhandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type DBMock struct{}

func NewDBMock() *DBMock {
	return &DBMock{}
}

func (db *DBMock) GetSiteDuration(name string) (time.Duration, bool) {
	LoadingTime := time.Duration(15923367)

	return LoadingTime, true
}

func (db *DBMock) CreateSite(name string, duration time.Duration) int64 {
	return 1
}

func (db *DBMock) GetReportByDate(from, to time.Time) (*[]Report, int64) {
	return &[]Report{
		{Name: "yandex.ru", AvgDuration: 12345678},
		{Name: "ya.ru", AvgDuration: 12345158}}, 1

}

type CacheMock struct{}

func NewCacheMock() *CacheMock {
	return &CacheMock{}
}

func (c *CacheMock) Get(name string) (time.Duration, bool) {
	if name == "http://yandex.ru" {
		return 15923367, true
	}

	return 0, false
}

func (c *CacheMock) Set(searchUrlname string, duration time.Duration) {

}

func TestCheckSite(t *testing.T) {

	db := NewDBMock()
	cache := NewCacheMock()

	handler := NewSiteHandler(30, 30, cache, db, zap.NewNop())

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)

	req, err := http.NewRequest("GET", "http://localhost:8081/sites?search=http://yandex.ru", nil)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	ctx.Request = req
	handler.CheckSite(ctx)

	responseBody := make(map[string]int, 1)

	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		t.Error("Unable to unmarshal JSON")
	}

	expectedResponseBody := map[string]int{
		"duration": 15923367,
	}

	require.Equal(t, responseBody, expectedResponseBody)

}

func TestGetReport(t *testing.T) {

	db := NewDBMock()
	cache := NewCacheMock()

	handler := NewSiteHandler(30, 30, cache, db, zap.NewNop())

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)

	req, err := http.NewRequest("GET", "http://localhost:8081/sites/stats?from=2022-03-01T00:04:05Z&to=2022-03-02T23:04:05Z", nil)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	ctx.Request = req
	handler.GetReport(ctx)

	responseBody := new([]Report)

	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		t.Error("Unable to unmarshal JSON")
	}

	expectedResponseBody := &[]Report{
		{Name: "yandex.ru", AvgDuration: 12345678},
		{Name: "ya.ru", AvgDuration: 12345158}}

	require.Equal(t, responseBody, expectedResponseBody)
}
