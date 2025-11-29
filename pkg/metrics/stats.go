package metrics

import "math"

// IntStats holds simple aggregate statistics over integer data.
type IntStats struct {
	Count int
	Sum   int
	Mean  float64
	Max   int
	Std   float64
	CV    float64 // coefficient of variation = Std / Mean
}

// ComputeIntStats computes basic statistics over a slice of ints.
//
// It returns zeroed stats if xs is empty.
func ComputeIntStats(xs []int) IntStats {
	var st IntStats
	n := len(xs)
	if n == 0 {
		return st
	}
	st.Count = n

	// Sum and max
	max := xs[0]
	sum := 0
	for _, v := range xs {
		sum += v
		if v > max {
			max = v
		}
	}
	st.Sum = sum
	st.Max = max
	st.Mean = float64(sum) / float64(n)

	// Std dev
	if n > 1 {
		var sq float64
		for _, v := range xs {
			d := float64(v) - st.Mean
			sq += d * d
		}
		st.Std = math.Sqrt(sq / float64(n))
	}
	if st.Mean != 0 {
		st.CV = st.Std / st.Mean
	}
	return st
}
