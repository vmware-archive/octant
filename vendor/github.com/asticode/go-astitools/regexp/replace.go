package astiregexp

import "regexp"

// ReplaceAll replaces all matches from a source
func ReplaceAll(rgx *regexp.Regexp, src *[]byte, rpl []byte) {
	// Find all matches
	var start, end, delta, offset, i int
	var l = len(rpl)
	for _, indexes := range rgx.FindAllIndex(*src, -1) {
		// Update indexes
		start = indexes[0] + offset
		end = indexes[1] + offset
		delta = (end - start) - l
		offset -= delta

		// Update src length
		if delta < 0 {
			// Insert
			(*src) = append((*src)[:start], append(make([]byte, -delta), (*src)[start:]...)...)
		} else if delta > 0 {
			// Delete
			(*src) = append((*src)[:start], (*src)[start+delta:]...)
		}

		// Update src content
		for i = 0; i < l; i++ {
			(*src)[i+start] = rpl[i]
		}
	}
}
