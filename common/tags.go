package common

import "strings"

// check tags contain features from a groups of whitelists
// trim leading/trailing spaces from keys and values

// check if a tag list is empty or not

func BuildTags(tagList string) map[string][]string {
	conditions := make(map[string][]string)
	for _, group := range strings.Split(tagList, ",") {
		conditions[group] = strings.Split(group, "+")
	}
	return conditions
}
