package mredis

import (
	"sync"
	"github.com/garyburd/redigo/redis"
	"time"
	"strings"
	"github.com/2liang/mcache/modules/utils/setting"
)

type RedisOption struct {
	timeout			int
	readtimeout		int
	writetimeout	int
	db 				int
	mhosts 			string
	shosts			string
}

type BaseRedis struct {
	Mutex *sync.Mutex
	mredis *redis.Pool
	sredis []*redis.Pool
}

func (b *BaseRedis) InitRedis(option *RedisOption)  {

	b.mredis = &redis.Pool{
		MaxIdle:		5,
		MaxActive: 		500,
		IdleTimeout: 	30 * time.Second,
		Dial:           func() (redis.Conn, error) {
			c, err := redis.DialTimeout("tcp", option.mhosts, time.Duration(option.timeout) * time.Second, time.Duration(option.readtimeout) * time.Second, time.Duration(option.writetimeout * time.Second))
			if err != nil {
				return nil, err
			}

			if _, err := c.Do("SELECT", option.db); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow:    func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}

			_, err := c.Do("PING")
			return err
		},
	}

	shosts := strings.Split(option.shosts, ",")
	if len(shosts) < 1 {
		setting.Logger.Panic("The redis slave init error")
	}

	b.sredis = make([]*redis.Pool, 0)
	for i := 0; i < len(shosts); i++ {
		func(i int) {
			sredis_pool := &redis.Pool{
				MaxIdle:		5,
				MaxActive: 		500,
				IdleTimeout: 	30 * time.Second,
				Dial:           func() (redis.Conn, error) {
					c, err := redis.DialTimeout("tcp", shosts[i], time.Duration(option.timeout) * time.Second, time.Duration(option.readtimeout) * time.Second, time.Duration(option.writetimeout * time.Second))
					if err != nil {
						return nil, err
					}

					if _, err := c.Do("SELECT", option.db); err != nil {
						c.Close()
						return nil, err
					}
					return c, nil
				},
				TestOnBorrow:    func(c redis.Conn, t time.Time) error {
					if time.Since(t) < time.Minute {
						return nil
					}

					_, err := c.Do("PING")
					return err
				},
			}
			b.sredis = append(b.sredis, sredis_pool)
		}(i)
	}
}

func (b *BaseRedis) getMredis() redis.Conn {
	return b.mredis.Get()
}

func (b *BaseRedis) getSredis() redis.Conn {
	return b.sredis[0].Get()
}

func (b *BaseRedis) Get(key string) (r interface{}, err error) {
	conn := b.getSredis()
	return conn.Do("GET", key)
}

func (b *BaseRedis) Set(v ...interface{}) (r interface{}, err error) {

	conn := b.getMredis()
	return conn.Do("SET", v...)
}

func (b *BaseRedis) Exists(k string) (r interface{}, err error) {

	conn := b.getSredis()
	return conn.Do("EXISTS", k)
}

func (b *BaseRedis) Expire(v ...interface{}) (r interface{}, err error) {
	conn := b.getMredis()
	return conn.Do("EXPIRE", v...)
}

func (b *BaseRedis) Incrby(v ...interface{}) (r interface{}, err error) {
	conn := b.getMredis()
	return conn.Do("INCRBY", v...)
}

func (b *BaseRedis) Decrby(v ...interface{}) (r interface{}, err error) {

	conn := b.getMredis()
	return conn.Do("DECRBY", v...)
}

func (b *BaseRedis) RPush(v ...interface{}) (r interface{}, err error) {

	conn := b.getMredis()
	return conn.Do("RPUSH", v...)
}

func (b *BaseRedis) LPop(k string) (r interface{}, err error) {

	conn := b.getMredis()
	return conn.Do("LPOP", k)
}

func (b *BaseRedis) LLen(k string) (r interface{}, err error) {

	conn := b.getSredis()
	return conn.Do("LLEN", k)
}

func (b *BaseRedis) LRem(v ...interface{}) (r interface{}, err error) {

	conn := b.getMredis()
	return conn.Do("LRem", v...)
}