package limiter

import "context"

// MockStorage implementa a interface Storage para testes
type MockStorage struct {
	Counts map[string]int64
}

func NewMockStorage() *MockStorage {
	return &MockStorage{
		Counts: make(map[string]int64),
	}
}

func (m *MockStorage) Increment(ctx context.Context, key string) (int64, error) {
	m.Counts[key]++
	return m.Counts[key], nil
}

func (m *MockStorage) Get(ctx context.Context, key string) (int64, error) {
	return m.Counts[key], nil
}

func (m *MockStorage) Set(ctx context.Context, key string, value int64, expiration int) error {
	m.Counts[key] = value
	return nil
}

func (m *MockStorage) Delete(ctx context.Context, key string) error {
	delete(m.Counts, key)
	return nil
}

func (m *MockStorage) Close() error {
	return nil
}
