package check

import "github.com/go-redis/redis/v7"

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
	if config.Redis != nil {
		password := ""
		if config.Redis.Password != nil {
			password = *config.Redis.Password
		}
		return &RedisStore{
			redis.NewClient(&redis.Options{
				Addr:     config.Redis.Address,
				Password: password,
			}),
		}
	}
	return &MemoryStore{make(map[string]string)}
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
