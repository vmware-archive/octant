package conversion

import "fmt"

// PtrInt32ToString convert *int32 to string
func PtrInt32ToString(p *int32) string {
	var i int32
	if p != nil {
		i = *p
	}

	return fmt.Sprintf("%d", i)
}
