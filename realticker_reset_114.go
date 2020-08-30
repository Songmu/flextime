// +build !go1.15

package flextime

import (
	"log"
	"time"
)

func (t *realTicker) Reset(d time.Duration) {
	log.Println("can't call ticker.Reset before Go 1.14")
}
