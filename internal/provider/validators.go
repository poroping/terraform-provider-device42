package provider

import "regexp"

func validateRegExpVlanRange() *regexp.Regexp {
	r, _ := regexp.Compile(`^(?:[1-9]\d{0,2}|[1-3]\d{3}|40(?:[0-8]\d|9[0-4]))(?:[,-] *(?:[1-9]\d{0,2}|[1-3]\d{3}|40(?:[0-8]\d|9[0-4]))?)*$`)
	return r
}
