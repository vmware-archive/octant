package strings

// Contains returns true if a slice contains a string and false if not.
func Contains(s string, sl []string) bool {
	for i := range sl {
		if s == sl[i] {
			return true
		}
	}

	return false
}
