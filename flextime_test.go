package flextime_test

import (
	"testing"
	"time"

	"github.com/Songmu/flextime"
)

func TestRestore(t *testing.T) {
	flextime.Fix(baseDate)
	almostSameTime(t, flextime.Now(), baseDate)

	flextime.Restore()
	almostSameTime(t, flextime.Now(), time.Now())
}
