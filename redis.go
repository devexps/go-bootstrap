package bootstrap

import (
	"github.com/devexps/go-micro/v2/log"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"

	conf "github.com/devexps/go-bootstrap/gen/api/go/conf/v1"
)

// NewRedisClient creates redis client
func NewRedisClient(conf *conf.Data) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:         conf.Redis.Addr,
		Password:     conf.Redis.Password,
		DB:           int(conf.Redis.Db),
		DialTimeout:  conf.Redis.DialTimeout.AsDuration(),
		WriteTimeout: conf.Redis.WriteTimeout.AsDuration(),
		ReadTimeout:  conf.Redis.ReadTimeout.AsDuration(),
	})
	if rdb == nil {
		log.Fatalf("failed opening connection to redis")
		return nil
	}
	// open tracing instrumentation.
	if conf.GetRedis().GetEnableTracing() {
		if err := redisotel.InstrumentTracing(rdb); err != nil {
			log.Fatalf("failed open tracing: %s", err.Error())
			panic(err)
		}
	}

	// open metrics instrumentation.
	if conf.GetRedis().GetEnableMetrics() {
		if err := redisotel.InstrumentMetrics(rdb); err != nil {
			log.Fatalf("failed open metrics: %s", err.Error())
			panic(err)
		}
	}
	return rdb
}
