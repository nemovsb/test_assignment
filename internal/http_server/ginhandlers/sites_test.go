package ginhandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"test_assignment/internal/configuration/di"
	"test_assignment/internal/storage"
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

func (db *DBMock) GetSiteByName(name string) (*storage.Sites, int64) {
	created, err := time.Parse("2006-01-02T15:04:05Z07:00", "2022-03-01T00:04:05Z")
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return nil, 0
	}
	var site storage.Sites

	site.ID = 1
	site.CreatedAt = created
	site.UpdatedAt = created
	site.Name = "http://yandex.ru"
	site.LoadingTime = 15923367

	return &site, 1
}

func (db *DBMock) CreateSite(name string, duration time.Duration) int64 {
	return 1
}

func (db *DBMock) GetReportByDate(from, to time.Time) (*[]storage.Report, int64) {
	return &[]storage.Report{
		{Name: "yandex.ru", AvgDuration: 12345678},
		{Name: "ya.ru", AvgDuration: 12345158}}, 1

}

func TestCheckSite(t *testing.T) {

	config := &di.ConfigApp{}
	config.TTL = 30000000000000000
	config.Timeout = 30

	db := NewDBMock()

	handler := NewSiteHandler(config, db, zap.NewNop())

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
	config := &di.ConfigApp{}
	config.TTL = 30
	config.Timeout = 30

	db := NewDBMock()

	handler := NewSiteHandler(config, db, zap.NewNop())

	rr := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(rr)

	req, err := http.NewRequest("GET", "http://localhost:8081/sites/stats?from=2022-03-01T00:04:05Z&to=2022-03-02T23:04:05Z", nil)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	ctx.Request = req
	handler.GetReport(ctx)

	responseBody := new([]storage.Report)

	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		t.Error("Unable to unmarshal JSON")
	}

	expectedResponseBody := &[]storage.Report{
		{Name: "yandex.ru", AvgDuration: 12345678},
		{Name: "ya.ru", AvgDuration: 12345158}}

	require.Equal(t, responseBody, expectedResponseBody)
}
