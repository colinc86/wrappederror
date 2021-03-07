package wrappederror

import (
	"fmt"
	"hash/fnv"
	"sync"
)

// A hash map of errors.
type errorMap struct {
	hashMap *sync.Map
}

// Initializers

// newErrorMap creates and returns a new error map.
func newErrorMap() *errorMap {
	return &errorMap{
		hashMap: new(sync.Map),
	}
}

// Non-exported methods

// similarErrors returns the number of similar errors.
func (m errorMap) similarErrors(err error) int {
	hash := string(m.hashError(err))
	if v, ok := m.hashMap.Load(hash); ok {
		return v.(int)
	}
	return 0
}

// addError adds an error to the map.
func (m *errorMap) addError(err error) {
	hash := string(m.hashError(err))
	if v, ok := m.hashMap.Load(hash); ok {
		m.hashMap.Store(hash, v.(int)+1)
	} else {
		m.hashMap.Store(hash, 1)
	}
}

// hashError hashes an error.
func (m errorMap) hashError(err error) []byte {
	s := fmt.Sprintf("%+v", err)
	h := fnv.New128a()
	h.Write([]byte(s))
	return h.Sum(nil)
}
