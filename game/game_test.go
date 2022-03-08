package game

import (
	"testing"
	"time"
)

func BenchmarkManager_Update(b *testing.B) {
	b.StopTimer()
	m := NewManager()
	m.InitGame()

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Update(time.Millisecond.Seconds())
	}
}
