package cfg

import (
	server "test_assignment/internal/http_server"
	"test_assignment/internal/storage"
)

func GetHTTPServerConfig(config *ConfigApp) server.ServerConfig {
	return server.ServerConfig{
		Port:    config.HttpServer.Port,
		Timeout: config.HttpServer.Timeout,
		TTL:     config.HttpServer.TTL,
	}
}

func GetDBConfig(config *ConfigApp) storage.DBConfig {
	return storage.DBConfig{
		Host:     config.DataBase.Host,
		Port:     config.DataBase.Port,
		User:     config.DataBase.User,
		Password: config.DataBase.Password,
		DBName:   config.DataBase.DBName,
	}
}

func GetCacheConfig(config *ConfigApp) storage.RedisConfig {
	return storage.RedisConfig{
		Host: config.Cache.Host,
		Port: config.Cache.Port,
	}
}
