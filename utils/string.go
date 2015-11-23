package utils

import (
	"strings"
)

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		return ""
	}
	str = string([]byte(str)[n+len(start):])
	m := strings.Index(str, end)
	if m == -1 {
		return ""
	}
	str = string([]byte(str)[:m])
	return str
}

func GetAfterNStr(str, start string, m int) string {
	n := strings.Index(str, start)
	if n == -1 || n+len(start)+m > len(str) {
		return ""
	}
	str = string([]byte(str)[n+len(start) : n+len(start)+m])
	return str
}

func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}
