package flextime

import "time"

// Fix switches backend Clock and fixes the current time to the specified time. This differs
// from `Set` in that the time does not change even if the time elapses during the test.
// However, when Sleep is called, the virtual time is passed without actually pausing.
// It returns a restore func and we can restore time behavior by calling it after the test.
func Fix(t time.Time) (restore func()) {
	return Switch(newFixedClock(t))
}

// Set switches backend Clock and sets the current time to the specified time.
// It internally holds the offset and is affected by the passage of time during the test.
// When Sleep is called, the virtual time is passed without actually pausing.
// It returns a restore func and we can restore time behavior by calling it after the test.
func Set(t time.Time) (restore func()) {
	return Switch(newOffsetClock(t))
}
