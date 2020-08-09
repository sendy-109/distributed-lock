package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sendy-109/distributed-lock/rds"
	"github.com/sendy-109/distributed-lock/rlock"
	"time"
)

func main(){

	//redis链接
	option := &redis.RingOptions{
		Addrs: make(map[string]string, 0),
		PoolSize:       10,
		PoolTimeout:    10,
		//Password:       "",
	}
	option.Addrs["0"] = "192.168.0.104:6379"
	rds, err := rds.New(option)
	if err != nil {
		panic(err)
	}
	defer rds.Close()

	//创建分布式锁  p
	l := rlock.NewLock("test", 15*time.Second, rds)
	err = l.Lock()
	if err != nil {
		panic(err)
	}

	t, err := l.GetTtl()
	if err != nil {
		panic(err)
	}
	fmt.Println(fmt.Sprintf("new lock sucess: %s, ttl:%d ms", l.GetKey(), t/time.Millisecond))

	time.Sleep(5*time.Second)
	//延长锁时间
	err = l.Expire()
	if err != nil {
		panic(err)
	}

	t, err = l.GetTtl()
	if err != nil {
		panic(err)
	}
	fmt.Println("Expire sucess "+ fmt.Sprintf("%d ms",t/time.Millisecond))

	//释放锁
	err = l.UnLock()
	if err != nil {
		panic(err)
	}

	fmt.Println("unlock sucess ")

}
