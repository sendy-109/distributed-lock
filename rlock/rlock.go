package rlock

import (
	"errors"
	"github.com/sendy-109/distributed-lock/rds"
	"time"
)

const(
	//重试间隔
	CONST_RETRY_INTER = 100*time.Millisecond
)

type  RdsLock interface {
	Lock() error
	Expire() error
	UnLock() error
	GetKey() string
	GetTtl() (time.Duration, error)
}

type locker struct {
	ttl              time.Duration     //过期时间
	retryInter       time.Duration     //重试间隔
	key              string            //锁标识
	rds rds.RdsClient
}

//创建锁
func NewLock(key string, ttl time.Duration, rds rds.RdsClient) RdsLock{
	l := locker{ttl:ttl, key:key, rds:rds, retryInter: CONST_RETRY_INTER}
	return &l
}

//获取锁
func (l *locker)Lock() error{

	var timer *time.Timer
	for deadline := time.Now().Add(l.ttl); time.Now().Before(deadline); {
		sucess, err := l.rds.SetNX(l.key, "redis lock", l.ttl).Result()
		if err != nil {
			return err
		} else if sucess{
			return nil
		}

		if timer == nil {
			timer = time.NewTimer(l.retryInter)
		} else {
			timer.Reset(l.retryInter)
		}

		select {
		case <-timer.C:
			continue
		}
	}
	return errors.New("lock time out")
}

//释放锁
func (l *locker)UnLock() error {
	num, err := l.rds.Del(l.key).Result()
	if err != nil {
		return err
	} else if (num != 1){
		return errors.New("UnLock err, inval key: "+ l.key)
	}
	return nil
}

//延长一个ttl时间
func (l *locker)Expire() error {
	b, err := l.rds.Expire(l.key, l.ttl).Result()
	if err != nil {
		return err
	}else if !b {
		return errors.New("Expire err, inval key: " + l.key)
	}
	return nil
}

//获取锁标识
func (l *locker)GetKey() string {
	return l.key
}

//获取过期时间
func (l *locker)GetTtl() (time.Duration, error){
	return  l.rds.PTTL(l.key).Result()
}
