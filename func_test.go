package flextime_test

import (
	"testing"
	"time"

	"github.com/Songmu/flextime"
)

var baseDate = time.Date(2080, time.June, 5, 22, 10, 10, 0, time.UTC)

func runTests(t *testing.T, fn func(t time.Time) func()) {
	t.Run("Tick", func(t *testing.T) {
		restore := fn(baseDate)
		defer restore()
		expect := baseDate

		interval := 13 * time.Second
		ch := flextime.Tick(interval)
		expect = expect.Add(interval)
		almostSameTime(t, <-ch, expect)

		expect = expect.Add(interval)
		almostSameTime(t, <-ch, expect)

		if flextime.Tick(-1) != nil {
			t.Errorf("Tick with negative value should return nil but not")
		}
	})

	restore := fn(baseDate)
	defer restore()

	expect := baseDate
	t.Run("Now", func(t *testing.T) {
		almostSameTime(t, flextime.Now(), expect)
	})

	sleep := 2 * time.Second
	t.Run("Sleep", func(t *testing.T) {
		flextime.Sleep(sleep)
		almostSameTime(t, flextime.Now(), expect.Add(sleep))
	})

	t.Run("Since", func(t *testing.T) {
		since := flextime.Since(expect)
		almostSameDuration(t, since, sleep)
	})

	t.Run("Until", func(t *testing.T) {
		until := flextime.Until(expect)
		almostSameDuration(t, -until, sleep)
	})

	expect = expect.Add(sleep)

	t.Run("After", func(t *testing.T) {
		after := 3 * time.Second
		got := <-flextime.After(after)
		expect = expect.Add(after)
		almostSameTime(t, got, expect)
	})

	t.Run("AfterFunc", func(t *testing.T) {
		after := 5 * time.Second
		var (
			fired bool
			done  = make(chan struct{})
		)
		ti := flextime.AfterFunc(after, func() {
			fired = true
			done <- struct{}{}
		})
		expect = expect.Add(after)
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
			after := 7 * time.Second
			ti.Reset(after)
			expect = expect.Add(after)
			<-done
			if !fired {
				t.Errorf("AfterFunc not fired")
			}
		})

		t.Run("Stop", func(t *testing.T) {
			time.Sleep(10 * time.Millisecond)
			if ti.Stop() {
				t.Errorf("Timer should be stopped but active")
			}
		})
	})

	t.Run("NewTimer", func(t *testing.T) {
		after := 4 * time.Second
		ti := flextime.NewTimer(after)
		expect = expect.Add(after)
		got := <-ti.C
		almostSameTime(t, got, expect)

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
			after := 6 * time.Second
			ti.Reset(after)
			expect = expect.Add(after)
			got := <-ti.C
			almostSameTime(t, got, expect)

			// A negative or zero duration return immediately
			ti.Reset(-after)
			got = <-ti.C
			almostSameTime(t, got, expect)
		})

		t.Run("Stop", func(t *testing.T) {
			time.Sleep(10 * time.Millisecond)
			if ti.Stop() {
				t.Errorf("Timer should be stopped but active")
			}
		})
	})

	t.Run("NewTicker", func(t *testing.T) {
		interval := 11 * time.Second
		ti := flextime.NewTicker(interval)
		expect = expect.Add(interval)
		almostSameTime(t, <-ti.C, expect)
		ti.Stop()
		select {
		case <-ti.C:
		default:
		}
	})
}

func almostSameDuration(t *testing.T, g, e time.Duration) {
	t.Helper()
	if int(g.Seconds()) != int(e.Seconds()) {
		t.Errorf("got: %s, expect: %s", g, e)
	}
}

func almostSame(g, e time.Time) bool {
	return g.Unix() != e.Unix() || g.Location().String() != e.Location().String()
}

func almostSameTime(t *testing.T, g, e time.Time) {
	t.Helper()
	if almostSame(g, e) {
		t.Errorf("got: %s, expect: %s", g, e)
	}
}

func TestFix(t *testing.T) {
	runTests(t, flextime.Fix)
}

func TestSet(t *testing.T) {
	runTests(t, flextime.Set)
}

func TestFix_NewTicker_withSleep(t *testing.T) {
	restore := flextime.Fix(baseDate)
	defer restore()

	expect := baseDate
	interval := 11 * time.Second
	ti := flextime.NewTicker(interval)
	sleep := 25 * time.Second
	flextime.Sleep(sleep)
	almostSameTime(t, <-ti.C, expect.Add(interval))
	expect = expect.Add(3 * interval) // 11x3 keep base time
	expect2 := expect.Add(interval)
	got := <-ti.C
	if !almostSame(got, expect) && !almostSame(got, expect2) {
		t.Errorf("got: %s, expect: %s", got, expect)
	}
}

func TestFix_fix(t *testing.T) {
	restore := flextime.Fix(baseDate)
	defer restore()

	time.Sleep(10 * time.Microsecond)
	if !flextime.Now().Equal(baseDate) {
		t.Errorf("time doesn't fixed")
	}
}

func TestSet_slide(t *testing.T) {
	restore := flextime.Set(baseDate)
	defer restore()

	offset := 10 * time.Microsecond
	time.Sleep(offset)
	if flextime.Now().Sub(baseDate.Add(offset)) <= 0 {
		t.Errorf("time doesn't slide")
	}
}
