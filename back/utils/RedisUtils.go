package utils

import (
	"github.com/go-redis/redis"
	"t/back/config"
)

// 声明一个全局的rdb变量
var Rdb *redis.Client = nil

// 初始化连接
func GetRedisClient() *redis.Client{
	if Rdb != nil{
		return Rdb
	}
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.ProjectConfig.RedisUrl,
		DB:       config.ProjectConfig.RedisDB,  // use default DB
	})
	return Rdb
}
