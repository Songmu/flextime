package flextime_test

import (
	"testing"
	"time"

	"github.com/Songmu/flextime"
)

func TestFakeTicker_Reset(t *testing.T) {
	now := time.Now()
	restore := flextime.Fix(now)
	defer restore()
	ti := flextime.NewTicker(time.Second)
	ti.Reset(10 * time.Millisecond)
	expect := now.Add(1*time.Second + 10*time.Millisecond)
	almostSameTime(t, <-ti.C, expect)
}

func TestFakeTicker_Reset_panic(t *testing.T) {
	restore := flextime.Fix(time.Now())
	expect := "non-positive interval for NewTicker"
	defer func() {
		err := recover()
		if err.(error).Error() != expect {
			t.Errorf("got %v\nwant %s", err, expect)
		}
		restore()
	}()
	flextime.NewTicker(-1)
}
