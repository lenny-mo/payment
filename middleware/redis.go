package middleware

import (
	"fmt"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

const (
	Password = ""
	Host     = "redis6380"
	Port     = 6379
	DB       = 0
	PoolSize = 100
)

func Init() (err error) {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", Host, Port),
		Password: Password,
		DB:       DB,
		PoolSize: PoolSize,
		//Password: viper.GetString("redis.password"),
		//DB:       viper.GetInt("redis.db"),
		//PoolSize: viper.GetInt("redis.poolsize"),
	})

	_, err = redisClient.Ping().Result()
	if err != nil {
		fmt.Printf("connect redis failed, err:%v\n", err)
		return
	}

	return
}

func Close() {
	redisClient.Close()
}

// TODO
func RedisSet(key string, data interface{}) bool {
	return true
}

// TODO
func RedisGet(id string) interface{} {
	return nil
}
