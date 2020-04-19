package check

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v7"
)

const (
	redisPrefixValue  = "v:"
	redisPrefixStatus = "s:"
)

var (
	ErrNotFound = errors.New("value not found")
)

type (
	Store interface {
		GetValue(jobID string) (string, error)
		SetValue(jobID string, value string) error
		GetStatus(jobID string) (*JobStatus, error)
		SetStatus(jobID string, status *JobStatus) error
	}
	MemoryStore struct {
		statuses map[string]JobStatus
		values   map[string]string
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
	return &MemoryStore{make(map[string]JobStatus), make(map[string]string)}
}

func NewRedisStore(client *redis.Client) *RedisStore {
	return &RedisStore{client}
}

func (s *MemoryStore) GetValue(jobID string) (string, error) {
	v, exist := s.values[jobID]
	if !exist {
		return "", ErrNotFound
	}
	return v, nil
}

func (s *MemoryStore) SetValue(jobID string, value string) error {
	s.values[jobID] = value
	return nil
}

func (s *MemoryStore) GetStatus(jobID string) (*JobStatus, error) {
	v, exist := s.statuses[jobID]
	if !exist {
		return nil, ErrNotFound
	}
	return &v, nil
}

func (s *MemoryStore) SetStatus(jobID string, status *JobStatus) error {
	fmt.Println(*status)
	s.statuses[jobID] = *status
	return nil
}

func (s *RedisStore) GetValue(jobID string) (string, error) {
	b, err := s.client.Get(redisPrefixValue + jobID).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrNotFound
		}
		return "", err
	}
	return string(b), nil
}

func (s *RedisStore) SetValue(jobID string, value string) error {
	return s.client.Set(redisPrefixValue+jobID, value, 0).Err()
}

func (s *RedisStore) GetStatus(jobID string) (*JobStatus, error) {
	b, err := s.client.Get(redisPrefixStatus + jobID).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}
		return nil, err
	}
	var status JobStatus
	if err := json.Unmarshal([]byte(b), &status); err != nil {
		return nil, err
	}
	return &status, nil
}

func (s *RedisStore) SetStatus(jobID string, status *JobStatus) error {
	b, err := json.Marshal(status)
	if err != nil {
		return err
	}
	return s.client.Set(redisPrefixStatus+jobID, string(b), 0).Err()
}
