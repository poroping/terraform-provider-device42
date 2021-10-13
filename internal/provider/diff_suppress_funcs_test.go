package provider

import (
	"fmt"
	"testing"
)

func TestIsFakeListEqual(t *testing.T) {
	var tests = []struct {
		s1, s2 string
		equal  bool
	}{
		{"L2-DOMAIN-ACI-L2,ZZXX,ZZXX-DMZ,TERRAFORMED", "ZZXX,ZZXX-DMZ,L2-DOMAIN-ACI-L2,TERRAFORMED", true},
		{"a,b,c", "a,c,b", true},
		{"a,b,c", "d,e,f", false},
		{"a,b,c", "", false},
		{"", "a,b,c", false},
		{"a,,b", "a,b", true},
	}

	for i, tt := range tests {
		testname := fmt.Sprintf("Testing comma separated 'list' equality, %v", i)
		t.Run(testname, func(t *testing.T) {
			ans := isFakeListEqual(tt.s1, tt.s2)
			if ans != tt.equal {
				t.Errorf("got %v, want %v", ans, tt.equal)
			}
		})
	}
}
