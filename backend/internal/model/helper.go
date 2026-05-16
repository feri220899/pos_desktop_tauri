package model

import "strings"

func parseJSONStringArray(s string) []string {
	var result []string
	if len(s) < 2 {
		return result
	}
	s = s[1 : len(s)-1]
	if s == "" {
		return result
	}
	start := -1
	for i, c := range s {
		if c == '"' && start == -1 {
			start = i + 1
		} else if c == '"' && start != -1 {
			result = append(result, s[start:i])
			start = -1
		}
	}
	return result
}

func sortDir(order string) string {
	if strings.ToLower(order) == "desc" {
		return "DESC"
	}
	return "ASC"
}
