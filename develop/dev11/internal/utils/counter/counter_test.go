package counter

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounter_Inc(t *testing.T) {
	goroutines := 16
	incsPerGoroutine := 10000

	cnt := New(0)

	wg := &sync.WaitGroup{}
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < incsPerGoroutine; i++ {
				cnt.Inc()
			}
			wg.Done()
		}()
	}

	wg.Wait()

	assert.Equal(t, int64(goroutines*incsPerGoroutine), cnt.Value())
}
