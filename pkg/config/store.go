package config

import (
	"net/url"

	"github.com/go-redis/redis/v7"
	"github.com/uphy/watch-web/pkg/watch"
	"github.com/uphy/watch-web/pkg/config/template"
)

type (
	StoreConfig struct {
		Redis *RedisConfig `json:"redis,omitempty"`
	}
	RedisConfig struct {
		Address   *template.TemplateString `json:"address"`
		Password  *template.TemplateString `json:"password"`
		RedisToGo *template.TemplateString `json:"redistogo"`
	}
)

func newStore(ctx *template.TemplateContext, config *StoreConfig) (watch.Store, error) {
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
				return nil, err
			}
		}
		if addr != "" {
			client := redis.NewClient(&redis.Options{
				Addr:     addr,
				Password: password,
			})
			return watch.NewRedisStore(client), nil
		}
	}
	return watch.NewMemoryStore(), nil
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
