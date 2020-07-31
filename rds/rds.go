package rds

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type RdsClient interface {
	SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Expire(key string, ttl time.Duration) *redis.BoolCmd
	Del(keys ...string) *redis.IntCmd
	PTTL(key string) *redis.DurationCmd
}

//初始化redis链接
func New(opt *redis.RingOptions) (*redis.Ring, error){

	client := redis.NewRing(opt)
	_,err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	fmt.Println("redis 初始化链接成功")
	return client, nil

}