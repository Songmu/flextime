package flextime_test

import (
	"math"
	"testing"
	"time"

	"github.com/Songmu/flextime"
)

func TestRealClock(t *testing.T) {
	t.Run("Tick", func(t *testing.T) {
		interval := 123 * time.Millisecond
		ch := flextime.Tick(interval)

		t1 := <-ch
		almostNow(t, t1)

		t2 := <-ch
		almostNow(t, t2)

		if g, e := t2.Sub(t1), 100*time.Millisecond; g < e {
			t.Errorf("t2.Sub(t1) less than %s: %s", e, g)
		}

		if flextime.Tick(-1) != nil {
			t.Errorf("Tick with negative value should return nil but not")
		}
	})

	base := time.Now()
	t.Run("Now", func(t *testing.T) {
		almostNow(t, flextime.Now())
	})

	sleep := 2 * time.Second
	t.Run("Sleep", func(t *testing.T) {
		flextime.Sleep(sleep)
		almostNow(t, flextime.Now())
	})

	t.Run("Since", func(t *testing.T) {
		since := flextime.Since(base)
		almostSameDuration(t, since, sleep)
	})

	t.Run("Until", func(t *testing.T) {
		until := flextime.Until(base)
		almostSameDuration(t, -until, sleep)
	})

	t.Run("After", func(t *testing.T) {
		got := <-flextime.After(200 * time.Microsecond)
		almostNow(t, got)
	})

	t.Run("AfterFunc", func(t *testing.T) {
		after := 300 * time.Microsecond
		var (
			fired bool
			done  = make(chan struct{})
		)
		ti := flextime.AfterFunc(after, func() {
			fired = true
			done <- struct{}{}
		})
		<-done
		if !fired {
			t.Errorf("AfterFunc not fired")
		}
		var blocked bool
		select {
		case <-ti.C:
		default:
			blocked = true
		}
		if !blocked {
			t.Errorf("Timer.C with AfterFunc should be blocked")
		}

		t.Run("Reset", func(t *testing.T) {
			fired = false
			after := 500 * time.Microsecond
			ti.Reset(after)
			<-done
			if !fired {
				t.Errorf("AfterFunc not fired")
			}
		})

		t.Run("Stop", func(t *testing.T) {
			if ti.Stop() {
				t.Errorf("Timer should be stopped but active")
			}
		})
	})

	t.Run("NewTimer", func(t *testing.T) {
		after := 700 * time.Microsecond
		ti := flextime.NewTimer(after)
		got := <-ti.C
		almostNow(t, got)

		var blocked bool
		select {
		case <-ti.C:
		default:
			blocked = true
		}
		if !blocked {
			t.Errorf("drained Timer.C should be blocked")
		}

		t.Run("Reset", func(t *testing.T) {
			after := 600 * time.Microsecond
			ti.Reset(after)
			got := <-ti.C
			almostNow(t, got)

			// A negative or zero duration return immediately
			ti.Reset(-after)
			got = <-ti.C
			almostNow(t, got)
		})

		t.Run("Stop", func(t *testing.T) {
			if ti.Stop() {
				t.Errorf("Timer should be stopped but active")
			}
		})
	})

	t.Run("NewTicker", func(t *testing.T) {
		interval := 10 * time.Millisecond
		ti := flextime.NewTicker(interval)
		ti.Reset(interval)
		almostNow(t, <-ti.C)
		ti.Stop()
		select {
		case <-ti.C:
			t.Errorf("ti.C should be blocked but not")
		default:
		}
	})

}

func almostNow(t *testing.T, g time.Time) {
	e := time.Now()
	t.Helper()
	if time.Duration(math.Abs(float64(g.Sub(e)))) > time.Millisecond ||
		g.Location().String() != e.Location().String() {

		t.Errorf("got: %s, expect: %s", g, e)
	}
}
