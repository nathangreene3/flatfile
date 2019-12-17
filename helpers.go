package flatfile

// max returns the maximum value.
func max(xs ...int) int {
	m := xs[0]
	for _, x := range xs {
		if m < x {
			m = x
		}
	}

	return m
}

// min returns the minimum value.
func min(xs ...int) int {
	m := xs[0]
	for _, x := range xs {
		if x < m {
			m = x
		}
	}

	return m
}
