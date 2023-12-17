package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

// const (
// 	Password = ""
// 	Host     = "redis6380"
// 	Port     = 6379
// 	DB       = 0
// 	PoolSize = 100
// )

func StartRDB(host string, port int64, db, poolSize int, password string) error {
	// 使用传入的参数初始化 Redis 连接配置
	redisAddr := fmt.Sprintf("%s:%d", host, port)
	options := &redis.Options{
		Addr:     redisAddr,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	}

	// 创建 Redis 客户端
	redisClient = redis.NewClient(options)

	// 使用 Ping() 方法测试连接，确保连接成功
	_, err := redisClient.Ping().Result()
	if err != nil {
		fmt.Printf("connect redis failed, err:%v\n", err)
		return err
	} else {
		fmt.Println("connect redis success")
	}

	return nil
}

func Close() {
	redisClient.Close()
}

func RedisSet(key string, data interface{}) bool {
	// 检查 Redis 客户端是否已初始化
	if redisClient == nil {
		fmt.Println("Redis client is not initialized")
		return false
	}

	// 数据序列化为bytes
	datastr, err := json.Marshal(data)
	if err != nil {
		fmt.Println("error during json marshal")
		return false
	}
	// 将数据存入 Redis 中
	err = redisClient.Set(key, datastr, 24*time.Hour).Err()
	if err != nil {
		fmt.Println("Error setting data in Redis:", err)
		return false
	}

	return true
}

func RedisGet(id string) interface{} {
	// 检查 Redis 客户端是否已初始化
	if redisClient == nil {
		fmt.Println("Redis client is not initialized")
		return nil
	}

	// 从 Redis 中获取数据
	val, err := redisClient.Get(id).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Printf("Key %v does not exist in Redis", id)
			return nil
		}
		fmt.Println("Error getting data from Redis:", err)
		return nil
	}

	return val
}

func RedisDelete(key string) bool {
	// 检查 Redis 客户端是否已初始化
	if redisClient == nil {
		fmt.Println("Redis client is not initialized")
		return false
	}

	_, err := redisClient.Del(key).Result()
	if err != nil {
		fmt.Println("Delete Redis key failed", err.Error())
		return false
	}

	return true
}
