package httpserver

import (
	"net/url"
	"strconv"
	"strings"
)

func qGet(q url.Values, keys ...string) string {
	for _, k := range keys {
		if v := strings.TrimSpace(q.Get(k)); v != "" {
			return v
		}
	}
	return ""
}

func parseInt(q url.Values, keys ...string) int {
	v := qGet(q, keys...)
	if v == "" {
		return 0
	}
	n, _ := strconv.Atoi(v)
	return n
}

func parseIntPtr(q url.Values, keys ...string) *int {
	v := qGet(q, keys...)
	if v == "" {
		return nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return nil
	}
	return &n
}

func parseFloatPtr(q url.Values, keys ...string) *float64 {
	v := qGet(q, keys...)
	if v == "" {
		return nil
	}
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return nil
	}
	return &f
}

func parseBool(q url.Values, keys ...string) bool {
	v := strings.ToLower(qGet(q, keys...))
	return v == "true" || v == "1"
}
