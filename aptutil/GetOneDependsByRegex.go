package aptutil

import (
	"fmt"
	"regexp"
)

func GetOneDependsByRegex(
	vars map[string]string,
	regexpFilter string,
) string {
	deps := vars["Depends"]
	r := regexp.MustCompile(regexpFilter)
	matches := r.FindAllStringSubmatch(deps, -1)
	if len(matches) != 1 {
		panic(fmt.Errorf("invalid matches %#v for %s", matches, deps))
	}
	return matches[0][1]
}
