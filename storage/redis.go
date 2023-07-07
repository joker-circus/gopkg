package storage

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type RedisPool interface {
	Do(args ...interface{}) *redis.Cmd
	SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Ping() *redis.StatusCmd
	Eval(script string, keys []string, args ...interface{}) *redis.Cmd
	Exists(keys ...string) *redis.IntCmd
	Expire(key string, expiration time.Duration) *redis.BoolCmd
	ExpireAt(key string, tm time.Time) *redis.BoolCmd
	HSet(key, field string, value interface{}) *redis.BoolCmd
	Pipeline() redis.Pipeliner
	ZCard(key string) *redis.IntCmd
	ZRange(key string, start, stop int64) *redis.StringSliceCmd
	ZRevRange(key string, start, stop int64) *redis.StringSliceCmd
	Close() error
	Scan(cursor uint64, match string, count int64) *redis.ScanCmd
	Get(key string) *redis.StringCmd
	Del(keys ...string) *redis.IntCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Watch(fn func(*redis.Tx) error, keys ...string) error
	TTL(key string) *redis.DurationCmd
}

func NewRedisPool(addr, password string, useCluster bool) (RedisPool, error) {
	var redisPool RedisPool
	if useCluster {
		redisPool = redisCluster(addr, password)
	} else {
		redisPool = redisClient(addr, password)
	}

	_, err := redisPool.Ping().Result()
	if err != nil {
		return nil, errors.Wrapf(err, "connect redis failed")
	}
	return redisPool, nil
}

func redisClient(addr, password string) RedisPool {
	redisClient := redis.NewClient(&redis.Options{
		Addr:         addr,     // use default Addr
		Password:     password, // no password set
		DB:           0,        // use default DB
		DialTimeout:  2 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})
	return interface{}(redisClient).(RedisPool)
}

func redisCluster(addr, password string) RedisPool {
	redisClusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        strings.Split(addr, ";"),
		Password:     password,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	})
	return interface{}(redisClusterClient).(RedisPool)
}

func GetAllKeys(key string, batchScan int64, redisPool RedisPool) ([]string, error) {
	var cursor uint64
	var result []string

	for {
		var keys []string
		var err error
		keys, cursor, err = redisPool.Scan(cursor, key, batchScan).Result()
		if err != nil {
			return nil, err
		}
		result = append(result, keys...)
		if cursor == 0 {
			break
		}
	}

	return result, nil
}

func GetLock(c RedisPool, lockName string, acquireTimeout, lockTimeOut time.Duration) (string, error) {
	value := "1"
	endTime := time.Now().Add(acquireTimeout).UnixNano()
	for time.Now().UnixNano() <= endTime {
		if success, err := c.SetNX(lockName, value, lockTimeOut).Result(); err != nil && err != redis.Nil {
			return "", err
		} else if success {
			return value, nil
		} else if c.TTL(lockName).Val() == -1 {
			c.Expire(lockName, lockTimeOut)
		}
		time.Sleep(time.Millisecond * 10)
	}
	return "", errors.New("timeout")
}

func ReleaseLock(c RedisPool, lockName, code string) bool {
	txf := func(tx *redis.Tx) error {
		if v, err := tx.Get(lockName).Result(); err != nil && err != redis.Nil {
			return err
		} else if v == code {
			_, err := tx.Pipelined(func(pipe redis.Pipeliner) error {
				pipe.Del(lockName)
				return nil
			})
			return err
		}
		return nil
	}

	for {
		if err := c.Watch(txf, lockName); err == nil {
			return true
		} else if err == redis.TxFailedErr {
			fmt.Println("watch key is modified, retry to release lock. err:", err.Error())
		} else {
			fmt.Println("err:", err.Error())
			return false
		}
	}
}
