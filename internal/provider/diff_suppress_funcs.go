package provider

import (
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func diffFakeListEqual(k, old, new string, d *schema.ResourceData) bool {
	if old == "" && new == "" {
		return true
	}
	return isFakeListEqual(old, new)
}

func isFakeListEqual(s1, s2 string) bool {
	l1 := strings.Split(s1, ",")
	l2 := strings.Split(s2, ",")
	sort.Strings(l1)
	sort.Strings(l2)

	return reflect.DeepEqual(l1, l2)
}
