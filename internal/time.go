package internal

import (
	"time"
	_ "unsafe"
)

func InTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

//go:noescape
//go:linkname walltime runtime.walltime
func walltime() (int64, int32)

// TODO refactor time.Now to walltime
// https://tpaschalis.github.io/golang-time-now/
func Now() uint64 {
	x, y := walltime()
	return uint64(x)*1e9 + uint64(y)
}
