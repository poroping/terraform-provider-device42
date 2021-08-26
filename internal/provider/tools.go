package provider

import (
	"fmt"
)

func intList(l []interface{}) []string {
	s := make([]string, len(l))
	for i, v := range l {
		s[i] = fmt.Sprint(v)
	}
	return s
}
