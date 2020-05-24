package store

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v7"
	"github.com/uphy/watch-web/pkg/domain"
)

const (
	redisPrefixValue  = "v:"
	redisPrefixStatus = "s:"
)

var (
	ErrNotFound = errors.New("value not found")
)

type (
	MemoryStore struct {
		jobStatuses map[string]domain.JobStatus
		jobValues   map[string]string
		values      map[string]string
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

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{make(map[string]domain.JobStatus), make(map[string]string), make(map[string]string)}
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client}
}

func (s *MemoryStore) SetTemp(key string, value string, expire time.Duration) error {
	s.values[key] = value
	return nil
}

func (s *MemoryStore) Get(key string) (string, error) {
	v, exist := s.values[key]
	if !exist {
		return "", ErrNotFound
	}
	return v, nil
}

func (s *MemoryStore) GetJobValue(jobID string) (string, error) {
	v, exist := s.jobValues[jobID]
	if !exist {
		return "", ErrNotFound
	}
	return v, nil
}

func (s *MemoryStore) SetJobValue(jobID string, value string) error {
	s.jobValues[jobID] = value
	return nil
}

func (s *MemoryStore) GetJobStatus(jobID string) (*domain.JobStatus, error) {
	v, exist := s.jobStatuses[jobID]
	if !exist {
		return nil, ErrNotFound
	}
	return &v, nil
}

func (s *MemoryStore) SetJobStatus(jobID string, status *domain.JobStatus) error {
	s.jobStatuses[jobID] = *status
	return nil
}

func (s *RedisStore) SetTemp(key string, value string, expire time.Duration) error {
	if err := s.client.Set(redisPrefixValue+key, value, 0).Err(); err != nil {
		return err
	}
	if expire > 0 {
		if err := s.client.Expire(key, expire).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (s *RedisStore) Get(key string) (string, error) {
	b, err := s.client.Get(redisPrefixValue + key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrNotFound
		}
		return "", err
	}
	return string(b), nil
}

func (s *RedisStore) GetJobValue(jobID string) (string, error) {
	b, err := s.client.Get(redisPrefixValue + jobID).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrNotFound
		}
		return "", err
	}
	return string(b), nil
}

func (s *RedisStore) SetJobValue(jobID string, value string) error {
	return s.client.Set(redisPrefixValue+jobID, value, 0).Err()
}

func (s *RedisStore) GetJobStatus(jobID string) (*domain.JobStatus, error) {
	b, err := s.client.Get(redisPrefixStatus + jobID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}
		return nil, err
	}
	var status domain.JobStatus
	if err := json.Unmarshal([]byte(b), &status); err != nil {
		return nil, err
	}
	return &status, nil
}

func (s *RedisStore) SetJobStatus(jobID string, status *domain.JobStatus) error {
	b, err := json.Marshal(status)
	if err != nil {
		return err
	}
	return s.client.Set(redisPrefixStatus+jobID, string(b), 0).Err()
}
