package check

import (
	"log"
	"net/url"

	"github.com/go-redis/redis/v7"
)

const redisPrefix = "v:"

type (
	Store interface {
		Get(name string) (*string, error)
		Set(name string, value string) error
	}
	MemoryStore struct {
		store map[string]string
	}
	RedisStore struct {
		client *redis.Client
	}
)

func newStore(config *StoreConfig) Store {
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
			return &RedisStore{
				redis.NewClient(&redis.Options{
					Addr:     addr,
					Password: password,
				}),
			}
		}
	}
	return &MemoryStore{make(map[string]string)}
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

func (s *MemoryStore) Get(name string) (*string, error) {
	v, exist := s.store[name]
	if !exist {
		return nil, nil
	}
	return &v, nil
}

func (s *MemoryStore) Set(name string, value string) error {
	s.store[name] = value
	return nil
}

func (s *RedisStore) Get(name string) (*string, error) {
	v, err := s.client.Get(redisPrefix + name).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, err
	}
	return &v, nil
}

func (s *RedisStore) Set(name string, value string) error {
	return s.client.Set(redisPrefix+name, value, 0).Err()
}
