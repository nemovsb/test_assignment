package storage

import (
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	db, err := gorm.Open(postgres.Open(url.String()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("can't create new db connection: %s", err)
	}

	db.AutoMigrate(&Sites{})

	return db, err
}
