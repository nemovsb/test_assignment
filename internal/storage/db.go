package storage

import (
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB interface {
	GetSiteByName(name string) (*Sites, int64)
	CreateSite(name string, duration time.Duration) int64
	GetReportByDate(from, to time.Time) (*[]Report, int64)
}

type PGDB struct {
	Connect *gorm.DB
}

func NewPGDB(config DBConfig) (*PGDB, error) {
	conn, err := NewDBConn(config)
	if err != nil {
		return nil, fmt.Errorf("postgres connect error: %s", err)
	}
	return &PGDB{
		Connect: conn,
	}, nil
}

func (db *PGDB) GetSiteByName(name string) (*Sites, int64) {
	var site Sites
	res := db.Connect.Where(`updated_at = (?)`, db.Connect.Table("sites").
		Select(`max(updated_at)`).
		Group(`name`).
		Where("name = ?", name)).
		Find(&site)
	return &site, res.RowsAffected
}

func (db *PGDB) CreateSite(name string, duration time.Duration) int64 {
	res := db.Connect.Create(&Sites{Name: name, LoadingTime: duration})
	return res.RowsAffected
}

func (db *PGDB) GetReportByDate(from, to time.Time) (*[]Report, int64) {
	report := []Report{}

	res := db.Connect.Table("sites").
		Select(`name, avg(loading_time)::bigint AS "avg_duration"`).
		Where("created_at >= ?", from).
		Where("created_at <= ?", to).
		Group("name").
		Find(&report)

	return &report, res.RowsAffected
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func NewDBConn(config DBConfig) (*gorm.DB, error) {

	//url := "postgres://myuser:secret@localhost:5432/mydb"

	var url strings.Builder
	url.WriteString(`postgres://`)
	url.WriteString(config.User)
	url.WriteString(`:`)
	url.WriteString(config.Password)
	url.WriteString(`@`)
	url.WriteString(config.Host)
	url.WriteString(`:`)
	url.WriteString(config.Port)
	url.WriteString(`/`)
	url.WriteString(config.DBName)

	fmt.Println(url.String()) //TODO

	db, err := gorm.Open(postgres.Open(url.String()), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now().Local().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("can't create new db connection: %s", err)
	}

	db.AutoMigrate(&Sites{})

	return db, err
}

// type DBMock struct{}

// func NewDBMock() *DBMock {
// 	return &DBMock{}
// }

// func (db *DBMock) GetSiteByName(name string) (Sites, int64) {
// 	created, err := time.Parse("2006-01-02T15:04:05Z07:00", "2022-03-01T00:04:05Z")
// 	if err != nil {
// 		fmt.Printf("error: %s\n", err)
// 		return nil, 0
// 	}
// 	var site Sites

// 	site.ID = 1
// 	site.CreatedAt = created
// 	site.UpdatedAt = created
// 	site.Name = "http://yandex.ru"
// 	site.LoadingTime = 15923367

// 	return &site, 1
// }

// func (db *DBMock) CreateSite(name string, duration time.Duration) int64 {
// 	return 1
// }

// func (db *DBMock) GetReportByDate(from, to time.Time) ([]Report, int64) {
// 	return []Report{
// 		{Name: "yandex.ru", AvgDuration: 12345678},
// 		{Name: "ya.ru", AvgDuration: 12345158}}, 1

// }
