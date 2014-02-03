package tb

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestNewBucket(t *testing.T) {
	t.Parallel()

	if b := NewBucket(10); b.tokens != 10 {
		t.Errorf("Wrong number of tokens. Want 10, Got %d", b.tokens)
	}
}

func TestBucket_Take_single(t *testing.T) {
	t.Parallel()

	b := NewBucket(10)

	ex := [...]int64{5, 5, 1, 1, 5, 4, 1, 0}
	for i := 0; i < len(ex)-1; i += 2 {
		if got, want := b.Take(ex[i]), ex[i+1]; got != want {
			t.Errorf("Want: %d, Got: %d", want, got)
		}
	}
}

func TestBucket_Take_multi(t *testing.T) {
	t.Parallel()

	b := NewBucket(10)

	exs := [2][]int64{{4, 4, 2, 2}, {2, 2, 1, 1}}
	for i := 0; i < 2; i++ {
		go func(i int) {
			for j := 0; j < len(exs[i])-1; j += 2 {
				if got, want := b.Take(exs[i][j]), exs[i][j+1]; got != want {
					t.Errorf("Want: %d, Got: %d", want, got)
				}
			}
		}(i)
	}
}

func TestBucket_Take_throughput(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	b := NewBucket(1000)
	atomic.StoreInt64(&b.tokens, 0)

	var out int64
	ts := time.Now()
	for time.Now().Before(ts.Add(1 * time.Second)) {
		out += b.Take(1)
	}

	// The time scheduler isn't as precise as we need so we need some tolerance
	thresholds := []int64{1000 - 50, 1000 + 50}
	if out < thresholds[0] || out > thresholds[1] {
		t.Errorf("Want %d to be within [%d, %d]", out, thresholds[0], thresholds[1])
	}
}

func BenchmarkBucket_Take_sequential(b *testing.B) {
	bucket := NewBucket(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bucket.Take(1)
	}
}
