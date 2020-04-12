package config

import (
	"log"
	"net/url"

	"github.com/go-redis/redis/v7"
	"github.com/uphy/watch-web/check"
)

type (
	StoreConfig struct {
		Redis *RedisConfig `json:"redis,omitempty"`
	}
	RedisConfig struct {
		Address   *TemplateString `json:"address"`
		Password  *TemplateString `json:"password"`
		RedisToGo *TemplateString `json:"redistogo"`
	}
)

func newStore(ctx *TemplateContext, config *StoreConfig) (check.Store, error) {
	if config != nil && config.Redis != nil {
		password := ""
		addr := ""
		if config.Redis.Address != nil {
			if config.Redis.Password != nil {
				p, err := config.Redis.Password.Evaluate(ctx)
				if err != nil {
					return nil, err
				}
				password = p
			}
			a, err := config.Redis.Address.Evaluate(ctx)
			if err != nil {
				return nil, err
			}
			addr = a
		} else if config.Redis.RedisToGo != nil {
			r, err := config.Redis.RedisToGo.Evaluate(ctx)
			if err != nil {
				return nil, err
			}
			addr, password, err = parseRedisToGoURL(r)
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}
		if addr != "" {
			client := redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: password,
			})
			return check.NewRedisStore(client), nil
		}
	}
	return &check.NullStore{}, nil
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
