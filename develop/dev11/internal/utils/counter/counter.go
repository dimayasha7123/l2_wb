package counter

import "sync"

// Counter struct
type Counter struct {
	mutex *sync.RWMutex
	value int64
}

// New constructor
func New(init int64) *Counter {
	return &Counter{
		mutex: &sync.RWMutex{},
		value: init,
	}
}

// Inc method
func (c *Counter) Inc() int64 {
	return c.Add(1)
}

// Add method
func (c *Counter) Add(value int64) int64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.value += value
	return c.value
}

// Value method
func (c *Counter) Value() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.value
}
