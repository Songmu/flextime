package flextime_test

import (
	"testing"
	"time"

	"github.com/Songmu/flextime"
)

func TestFakeTicker_Reset(t *testing.T) {
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
