package utils

import (
	"fmt"
	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var Rdb *redis.Client

// 初始化连接
func InitClient() (err error) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "192.168.3.67:6379",
		DB:       0,  // use default DB
	})
	//
	fmt.Println("init redis")
	return nil
}
