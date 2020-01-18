package flextime_test

import (
	"sync"
	"testing"
	"time"

	"github.com/Songmu/flextime"
)

type ns struct {
}

var _ flextime.NowSleeper = (*ns)(nil)

func (n *ns) Now() time.Time {
	return time.Now()
}

func (n *ns) Sleep(d time.Duration) {
	time.Sleep(d)
}

func TestTimer_Stop(t *testing.T) {
	restore := flextime.Switch(flextime.NewFakeClock(&ns{}))
	defer restore()

	var (
		ti    = flextime.NewTimer(time.Second)
		fired = time.Now().Add(500 * time.Millisecond)
		trial = 5
		wg    = sync.WaitGroup{}
		done  = make(chan struct{})
	)
	wg.Add(trial)
	for i := 0; i < trial; i++ {
		// Call Stop() several times at the same time
		time.AfterFunc(time.Until(fired), func() {
			defer wg.Done()
			ti.Stop()
		})
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Errorf("Timer.Stop() should not be blocked")
	}
}
