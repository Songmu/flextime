// +build go1.15

package flextime

import "time"

func (t *realTicker) Reset(d time.Duration) {
	t.t.Reset(d)
}
