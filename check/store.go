package check

import (
	"encoding/json"
	"log"
	"net/url"

	"github.com/go-redis/redis/v7"
)

const redisPrefix = "v:"

type (
	Store interface {
		GetJob(name string, job *Job) error
		SetJob(name string, job *Job) error
	}
	NullStore struct {
	}
	RedisStore struct {
		client *redis.Client
	}
	RedisJob struct {
		LastUpdatedSec int64   `json:"l"`
		Value          *string `json:"v"`
		Error          *string `json:"e"`
		Status         string  `json:"s"`
		Count          int     `json:"c"`
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
	return &NullStore{}
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

func (s *NullStore) GetJob(name string, job *Job) error {
	return nil
}

func (s *NullStore) SetJob(name string, job *Job) error {
	return nil
}

func (s *RedisStore) GetJob(name string, job *Job) error {
	b, err := s.client.Get(redisPrefix + name).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	var redisJob RedisJob
	if err := json.Unmarshal([]byte(b), &redisJob); err != nil {
		return err
	}
	job.Error = redisJob.Error
	job.Status = Status(redisJob.Status)
	job.Count = redisJob.Count
	job.Previous = redisJob.Value
	return nil
}

func (s *RedisStore) SetJob(name string, job *Job) error {
	b, err := json.Marshal(&RedisJob{
		Error:  job.Error,
		Status: string(job.Status),
		Count:  job.Count,
		Value:  job.Previous,
	})
	if err != nil {
		return err
	}
	return s.client.Set(redisPrefix+name, string(b), 0).Err()
}
