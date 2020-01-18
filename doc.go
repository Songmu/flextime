/*
Package flextime improves time testability by replacing the backend clock flexibly.

It has a set of following 9 functions similar to the standard time package, making it easy to migrate
from standard time package.

	now := flextime.Now()
	flextime.Sleep()
	d := flextime.Until(date)
	d := flextime.Since(date)
	<-flextime.After(5*time.Second)
	flextime.AfterFunc(5*time.Second, func() { fmt.Println("Done") })
	timer := flextime.NewTimer(10*time.Second)
	ticker := flextime.NewTicker(10*time.Second)
	ch := flextime.Tick(3*time.Second)

By default, it behaves the same as the standard time package, but allows us to change or fix
the current time by using `Fix` and `Set` function.

	func () { // Set time
		restore := flextime.Set(time.Date(2001, time.May, 1, 10, 10, 10, 0, time.UTC))
		defer restore()

		now = flextime.Now() // returned set time
	}()

	func () { // Fix time
		restore := flextime.Fix(time.Date(2001, time.May, 1, 10, 10, 10, 0, time.UTC))
		defer restore()

		now = flextime.Now() // returned fixed time
	}()

Also, we can replace the backend clock by implementing our own `Clock` interface and combining
it with the Switch function.

	restore := flextime.Switch(clock)
	defer restore()
*/
package flextime
