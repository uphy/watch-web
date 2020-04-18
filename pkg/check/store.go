package check

import (
	"encoding/json"
	"time"

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

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client}
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
	lastUpdated := time.Unix(redisJob.LastUpdatedSec, 0)
	job.Last = &lastUpdated
	return nil
}

func (s *RedisStore) SetJob(name string, job *Job) error {
	b, err := json.Marshal(&RedisJob{
		Error:          job.Error,
		Status:         string(job.Status),
		Count:          job.Count,
		Value:          job.Previous,
		LastUpdatedSec: job.Last.Local().Unix(),
	})
	if err != nil {
		return err
	}
	return s.client.Set(redisPrefix+name, string(b), 0).Err()
}
