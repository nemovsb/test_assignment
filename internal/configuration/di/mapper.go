package di

import (
	server "test_assignment/internal/http_server"
	"test_assignment/internal/storage"
)

func GetHTTPServerConfig(config *ConfigApp) server.ServerConfig {
	return server.ServerConfig{
		Port:     config.HttpServer.Port,
		RTimeout: config.HttpServer.RTimeout,
		WTimeout: config.HttpServer.WTimeout,
	}
}

func GetDBConfig(config *ConfigApp) storage.DBConfig {
return storage.DBConfig{
	Host:config.
    Port     string
    User     string
    Password string
    DB       string
}
}
