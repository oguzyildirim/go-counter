package internal

import (
	"testing"
	"time"
)

var timeTests = []struct {
	start, end, check time.Time
	want              bool
}{
	{time.Now(), time.Now().Add(60 * time.Second), time.Now().Add(30 * time.Second), true},    // in the time span
	{time.Now(), time.Now().Add(60 * time.Second), time.Now().Add(110 * time.Second), false},  // not in the time span
	{time.Now(), time.Now().Add(60 * time.Second), time.Now().Add(-110 * time.Second), false}, // not in the time span
}

func TestInTimeSpan(t *testing.T) {
	for _, a := range timeTests {
		if got := InTimeSpan(a.start, a.end, a.check); got != a.want {
			t.Fatalf("%t, want %t", got, a.want)
		}
	}
	t.Log(len(timeTests), "test cases")
}
