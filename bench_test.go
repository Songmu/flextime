package flextime_test

import (
	"sync"
	"testing"
	"time"

	"github.com/Songmu/flextime"
)

func BenchmarkFlextime_Now(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		flextime.Now()
	}
}

func BenchmarkStd_Now(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		time.Now()
	}
}

func BenchmarkFlextime_Now_concur(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		go func() {
			flextime.Now()
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkStd_Now_concur(b *testing.B) {
	var wg sync.WaitGroup
	wg.Add(b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		go func() {
			time.Now()
			wg.Done()
		}()
	}
	wg.Wait()
}
