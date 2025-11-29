package metrics

import "testing"

func TestComputeIntStatsBasic(t *testing.T) {
	xs := []int{1, 2, 3, 4}
	st := ComputeIntStats(xs)
	if st.Count != 4 {
		t.Fatalf("expected count=4, got %d", st.Count)
	}
	if st.Sum != 10 {
		t.Fatalf("expected sum=10, got %d", st.Sum)
	}
	if st.Max != 4 {
		t.Fatalf("expected max=4, got %d", st.Max)
	}
	if st.Mean != 2.5 {
		t.Fatalf("expected mean=2.5, got %f", st.Mean)
	}
}
