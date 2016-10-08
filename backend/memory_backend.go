package backend

import (
	"errors"
)

// MemStore is a key value store in memory.
type MemStore struct {
	Data map[string][]string
}

// Init initializes the data store
func (m *MemStore) Init() error {
	m.Data = make(map[string][]string)
	return nil
}

// Close "closes" the store.. really we don't do anything because RAM IS CHEAP, RIGHT?!
func (m *MemStore) Close() {
}

// Add inserts a key into the backend
func (m *MemStore) Add(user string, key string) (bool, error) {
	m.Data[user] = append(m.Data[user], key)
	return true, nil
}

// RM removes a key from the backend
func (m *MemStore) RM(user string, key string) (bool, error) {
	var idx = -1
	for i := 0; i < len(m.Data[user]); i++ {
		if m.Data[user][i] == key {
			idx = i
		}
	}

	if idx == -1 {
		return false, errors.New("No key found!")
	}

	m.Data[user] = append(m.Data[user][:idx], m.Data[user][idx+1:]...)

	return true, nil
}

// RMAll removes all of the keys for a given user
func (m *MemStore) RMAll(user string) (bool, error) {
	delete(m.Data, user)
	return true, nil
}

// Get queries the backend for 'user' and returns all the keys
func (m *MemStore) Get(user string) ([]string, error) {
	return m.Data[user], nil
}

// GetAll returns the entire datastructure
func (m *MemStore) GetAll() ([]string, error) {
	var s []string
	for i := range m.Data {
		for j := range m.Data[i] {
			s = append(s, m.Data[i][j])
		}
	}
	return s, nil
}

// GetKeyCount queries the number of keys a given user has
func (m *MemStore) GetKeyCount(user string) (int, error) {
	return len(m.Data[user]), nil
}

// GetCount queries the number of keys stored in a backend
func (m *MemStore) GetCount() (int, error) {
	var i int
	for j := range m.Data {
		for _ = range m.Data[j] {
			i++
		}
	}
	return i, nil
}
