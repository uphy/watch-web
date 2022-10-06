package store

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/uphy/watch-web/pkg/domain"
)

const (
	fileNameVars  = "vars"
	jobFilePrefix = "job-"
)

type (
	DirectoryStore struct {
		directory string
	}
	VariableMap map[string]Variable
	Variable    struct {
		Value  string    `json:"value"`
		Expire time.Time `json:"expire"`
	}
	JobFile struct {
		Value  string            `json:"value"`
		Status *domain.JobStatus `json:"status"`
	}
)

func NewDirectoryStore(directory string) (*DirectoryStore, error) {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err := os.MkdirAll(directory, 0700)
		if err != nil {
			return nil, fmt.Errorf("failed to make data directory: %w", err)
		}
	}
	return &DirectoryStore{directory}, nil
}

func (s *DirectoryStore) file(name string) string {
	return filepath.Join(s.directory, name+".json")
}

func (s *DirectoryStore) read(name string, value interface{}) error {
	file := s.file(name)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil
	}
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&value)
}

func (s *DirectoryStore) write(name string, value interface{}) error {
	f, err := os.Create(s.file(name))
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(value)
}

func (s *DirectoryStore) SetTemp(key string, value string, expire time.Duration) error {
	var vars VariableMap
	if err := s.read(fileNameVars, vars); err != nil {
		return err
	}
	if vars == nil {
		vars = make(VariableMap, 0)
	}
	vars[key] = Variable{Value: value, Expire: time.Now().Add(expire)}
	return s.write(fileNameVars, vars)
}

func (s *DirectoryStore) Get(key string) (string, error) {
	var vars VariableMap
	if err := s.read(fileNameVars, &vars); err != nil {
		return "", err
	}
	if vars == nil {
		return "", ErrNotFound
	}
	v, exist := vars[key]
	if !exist {
		return "", ErrNotFound
	}
	if time.Now().After(v.Expire) {
		return "", ErrNotFound
	}
	return v.Value, nil
}

func (s *DirectoryStore) readJob(jobID string) (*JobFile, error) {
	var file *JobFile
	err := s.read(jobFilePrefix+jobID, &file)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (s *DirectoryStore) writeJob(jobID string, file *JobFile) error {
	return s.write(jobFilePrefix+jobID, file)
}

func (s *DirectoryStore) GetJobValue(jobID string) (string, error) {
	f, err := s.readJob(jobID)
	if err != nil {
		return "", err
	}
	if f == nil {
		return "", ErrNotFound
	}
	return f.Value, nil
}

func (s *DirectoryStore) SetJobValue(jobID string, value string) error {
	f, err := s.readJob(jobID)
	if err != nil {
		return err
	}
	if f == nil {
		f = &JobFile{}
	}
	f.Value = value
	return s.writeJob(jobID, f)
}

func (s *DirectoryStore) GetJobStatus(jobID string) (*domain.JobStatus, error) {
	f, err := s.readJob(jobID)
	if err != nil {
		return nil, err
	}
	if f == nil {
		return nil, ErrNotFound
	}
	return f.Status, nil
}

func (s *DirectoryStore) SetJobStatus(jobID string, status *domain.JobStatus) error {
	f, err := s.readJob(jobID)
	if err != nil {
		return err
	}
	if f == nil {
		f = &JobFile{}
	}
	f.Status = status
	return s.writeJob(jobID, f)
}
