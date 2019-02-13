package printer

import (
	"fmt"
)

func ptrInt32ToString(p *int32) string {
	var i int32
	if p != nil {
		i = *p
	}

	return fmt.Sprintf("%d", i)
}
