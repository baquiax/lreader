package inmemory

import "sync"

type Storage struct {
	ForcedReadError  error
	ForcedWriteError error
	mutex            sync.RWMutex
	value            int64
}

func (s *Storage) Read() (int64, error) {
	s.mutex.Lock()
	value := s.value
	s.mutex.Unlock()

	return value, s.ForcedReadError
}

func (s *Storage) Write(value int64) error {
	s.mutex.Lock()
	s.value = value
	s.mutex.Unlock()

	return s.ForcedWriteError
}
