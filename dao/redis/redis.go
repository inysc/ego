package redis

import "github.com/go-redis/redis/v9"

var rds *redis.Client

func InitRedis(addr, pwd string, idx int) *redis.Client {
	rds = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       idx,
	})
	return rds
}

func Close() {
	rds.Close()
}
