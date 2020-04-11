package config

import (
	"log"
	"net/url"

	"github.com/go-redis/redis/v7"
	"github.com/uphy/watch-web/check"
)

func (c *Config) NewExecutor() (*check.Executor, error) {
	store := newStore(c.Store)
	e := check.NewExecutor(store)
	if c.InitialRun != nil {
		e.InitialRun = *c.InitialRun
	}
	for name, job := range c.Jobs {
		source, err := job.Source.Source()
		if err != nil {
			return nil, err
		}
		actions := []check.Action{}
		for _, a := range job.Actions {
			action, err := a.Action()
			if err != nil {
				return nil, err
			}
			actions = append(actions, action)
		}
		e.AddJob(name, job.Schedule, job.Label, job.Link, source, actions...)
	}
	return e, nil
}

func newStore(config *StoreConfig) check.Store {
	if config != nil && config.Redis != nil {
		password := ""
		addr := ""
		if config.Redis.Address != nil {
			if config.Redis.Password != nil {
				password = *config.Redis.Password
			}
			addr = *config.Redis.Address
		} else if config.Redis.RedisToGo != nil {
			var err error
			addr, password, err = parseRedisToGoURL(*config.Redis.RedisToGo)
			if err != nil {
				log.Println(err)
				return nil
			}
		}
		if addr != "" {
			client := redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: password,
			})
			return check.NewRedisStore(client)
		}
	}
	return &check.NullStore{}
}

func parseRedisToGoURL(redisToGo string) (addr string, password string, err error) {
	var redisInfo *url.URL
	redisInfo, err = url.Parse(redisToGo)
	if err != nil {
		return
	}

	addr = redisInfo.Host
	if redisInfo.User != nil {
		password, _ = redisInfo.User.Password()
	}
	return
}
