package store

import (
	"time"

	"github.com/uphy/watch-web/pkg/domain"
)

type (
	MemoryStore struct {
		jobStatuses map[string]domain.JobStatus
		jobValues   map[string]string
		values      map[string]string
	}
)

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{make(map[string]domain.JobStatus), make(map[string]string), make(map[string]string)}
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
