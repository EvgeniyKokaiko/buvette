package utils

import (
	"fmt"
	"strings"
)

func Contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func TrimSlice(s []string) []string {
	items := []string{}
	for _, v := range s {
		items = append(items, strings.TrimSpace(v))
	}
	return items
}

func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func IsUpperCased(value string) bool {
	if value == strings.ToUpper(value) {
		return true
	}
	return false
}

func LogError(err error) {
	if err != nil {
		fmt.Println("[BUVETTE]: ERROR!", err)
		return
	}
}
