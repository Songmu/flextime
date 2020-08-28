// +build !go1.15

package flextime

import "time"

func (t *realTicker) Reset(d time.Duration) {
	panic("can't call ticker.Reset before Go 1.14")
}
