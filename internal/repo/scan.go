package repo

import (
	"database/sql"
	"strconv"
	"strings"
)

func rowsToMaps(rows *sql.Rows) ([]map[string]interface{}, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	var out []map[string]interface{}
	for rows.Next() {
		raw := make([]interface{}, len(cols))
		ptr := make([]interface{}, len(cols))
		for i := range raw {
			ptr[i] = &raw[i]
		}
		if err := rows.Scan(ptr...); err != nil {
			return nil, err
		}
		m := make(map[string]interface{}, len(cols))
		for i, c := range cols {
			m[strings.ToLower(strings.TrimSpace(c))] = raw[i]
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func mapInt(m map[string]interface{}, keys ...string) int {
	v := pick(m, keys...)
	if v == nil {
		return 0
	}
	switch x := v.(type) {
	case int64:
		return int(x)
	case int32:
		return int(x)
	case int:
		return x
	case []byte:
		n, _ := strconv.Atoi(strings.TrimSpace(string(x)))
		return n
	case float64:
		return int(x)
	default:
		return 0
	}
}

func mapIntPtr(m map[string]interface{}, keys ...string) *int {
	v := pick(m, keys...)
	if v == nil {
		return nil
	}
	n := mapInt(m, keys...)
	return &n
}

func mapFloatPtr(m map[string]interface{}, keys ...string) *float64 {
	v := pick(m, keys...)
	if v == nil {
		return nil
	}
	f, ok := toFloat64(v)
	if !ok {
		return nil
	}
	return &f
}

func mapString(m map[string]interface{}, keys ...string) string {
	v := pick(m, keys...)
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	default:
		return ""
	}
}

func pick(m map[string]interface{}, keys ...string) interface{} {
	for _, k := range keys {
		if v, ok := m[strings.ToLower(k)]; ok && v != nil {
			return v
		}
	}
	return nil
}

// Exported aliases for use outside repo row helpers.
func MapInt(m map[string]interface{}, keys ...string) int       { return mapInt(m, keys...) }
func MapString(m map[string]interface{}, keys ...string) string { return mapString(m, keys...) }
func MapFloatPtr(m map[string]interface{}, keys ...string) *float64 {
	return mapFloatPtr(m, keys...)
}
func Pick(m map[string]interface{}, keys ...string) interface{} { return pick(m, keys...) }
func ToFloat64(v interface{}) (float64, bool)                     { return toFloat64(v) }

func toFloat64(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int64:
		return float64(x), true
	case int32:
		return float64(x), true
	case []byte:
		f, err := strconv.ParseFloat(strings.TrimSpace(string(x)), 64)
		return f, err == nil
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(x), 64)
		return f, err == nil
	default:
		return 0, false
	}
}
