package uconv

import (
	"fmt"
	"strconv"
	"time"

	"github.com/whosafe/uf/uerror"
)

// ToString 将任意值转换为字符串
func ToString(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case *string:
		if val == nil {
			return ""
		}
		return *val
	case []byte:
		return string(val)
	case fmt.Stringer:
		return val.String()
	case error:
		return val.Error()
	default:
		return fmt.Sprintf("%v", v)
	}
}

// ToBool 转换为 bool
func ToBool(v any) (bool, error) {
	if v == nil {
		return false, uerror.New("value is nil")
	}
	switch val := v.(type) {
	case bool:
		return val, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		i, _ := ToInt64(val)
		return i != 0, nil
	case string:
		return strconv.ParseBool(val)
	default:
		return strconv.ParseBool(ToString(val))
	}
}

// ToBoolDef 如果转换失败返回默认值
func ToBoolDef(v any, def bool) bool {
	if b, err := ToBool(v); err == nil {
		return b
	}
	return def
}

// ToInt 将任意值转换为 int
func ToInt(v any) (int, error) {
	if v == nil {
		return 0, uerror.New("value is nil")
	}
	switch val := v.(type) {
	case int:
		return val, nil
	case int8:
		return int(val), nil
	case int16:
		return int(val), nil
	case int32:
		return int(val), nil
	case int64:
		return int(val), nil
	case uint:
		return int(val), nil
	case uint8:
		return int(val), nil
	case uint16:
		return int(val), nil
	case uint32:
		return int(val), nil
	case uint64:
		return int(val), nil
	case float32:
		return int(val), nil
	case float64:
		return int(val), nil
	case string:
		if val == "" {
			return 0, nil
		}
		return strconv.Atoi(val)
	default:
		// 尝试转换为字符串再解析
		s := fmt.Sprintf("%v", v)
		return strconv.Atoi(s)
	}
}

// ToIntDef 如果转换失败返回默认值
func ToIntDef(v any, def int) int {
	if i, err := ToInt(v); err == nil {
		return i
	}
	return def
}

// MustToInt 转换失败panic
func MustToInt(v any) int {
	i, err := ToInt(v)
	if err != nil {
		panic(err)
	}
	return i
}

// ToInt64 将任意值转换为 int64
func ToInt64(v any) (int64, error) {
	if v == nil {
		return 0, uerror.New("value is nil")
	}
	switch val := v.(type) {
	case int:
		return int64(val), nil
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return val, nil
	case uint:
		return int64(val), nil
	case uint8:
		return int64(val), nil
	case uint16:
		return int64(val), nil
	case uint32:
		return int64(val), nil
	case uint64:
		return int64(val), nil
	case float32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case string:
		if val == "" {
			return 0, nil
		}
		return strconv.ParseInt(val, 10, 64)
	default:
		s := fmt.Sprintf("%v", v)
		return strconv.ParseInt(s, 10, 64)
	}
}

// ToInt64Def 如果转换失败返回默认值
func ToInt64Def(v any, def int64) int64 {
	if i, err := ToInt64(v); err == nil {
		return i
	}
	return def
}

// ToUint64 将任意值转换为 uint64
func ToUint64(v any) (uint64, error) {
	if v == nil {
		return 0, uerror.New("value is nil")
	}
	switch val := v.(type) {
	case uint64:
		return val, nil
	case int, int8, int16, int32, int64:
		i, _ := ToInt64(val)
		if i < 0 {
			return 0, uerror.New("negative value cannot convert to uint64")
		}
		return uint64(i), nil
	case uint, uint8, uint16, uint32:
		i, _ := ToInt64(val) // safe since internal call to ToInt64 handles uints correctly by casting logic, but wait, ToInt64 handles uints.
		// ToInt64 implementation handles uints by casting to int64. Large uint64 might overflow int64.
		// Let's implement ToUint64 more carefully.
		return uint64(i), nil
	case float32, float64:
		f, _ := ToFloat64(val)
		if f < 0 {
			return 0, uerror.New("negative value cannot convert to uint64")
		}
		return uint64(f), nil
	case string:
		if val == "" {
			return 0, nil
		}
		return strconv.ParseUint(val, 10, 64)
	default:
		s := ToString(v)
		return strconv.ParseUint(s, 10, 64)
	}
}

// ToUint64Def
func ToUint64Def(v any, def uint64) uint64 {
	if i, err := ToUint64(v); err == nil {
		return i
	}
	return def
}

// ToFloat64 将任意值转换为 float64
func ToFloat64(v any) (float64, error) {
	if v == nil {
		return 0, uerror.New("value is nil")
	}
	switch val := v.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case uint:
		return float64(val), nil
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case string:
		if val == "" {
			return 0, nil
		}
		return strconv.ParseFloat(val, 64)
	default:
		s := ToString(v)
		return strconv.ParseFloat(s, 64)
	}
}

// ToFloat64Def
func ToFloat64Def(v any, def float64) float64 {
	if f, err := ToFloat64(v); err == nil {
		return f
	}
	return def
}

// ToDuration 将任意值转换为 time.Duration
// 支持 int(秒/纳秒? usually ns in Go), string("10s")
func ToDuration(v any) (time.Duration, error) {
	if v == nil {
		return 0, uerror.New("value is nil")
	}
	switch val := v.(type) {
	case time.Duration:
		return val, nil
	case int64:
		return time.Duration(val), nil
	case int:
		return time.Duration(val), nil
	case string:
		return time.ParseDuration(val)
	default:
		// Try to parse string
		return time.ParseDuration(ToString(v))
	}
}

// ToDurationDef
func ToDurationDef(v any, def time.Duration) time.Duration {
	if d, err := ToDuration(v); err == nil {
		return d
	}
	return def
}

// ToTime 将任意值转换为 time.Time
// 支持类型:
// - time.Time: 直接返回
// - int64/int: Unix时间戳(秒)
// - string: 尝试多种常见格式解析
//   - RFC3339: "2006-01-02T15:04:05Z07:00"
//   - DateTime: "2006-01-02 15:04:05"
//   - Date: "2006-01-02"
func ToTime(v any) (time.Time, error) {
	if v == nil {
		return time.Time{}, uerror.New("value is nil")
	}

	switch val := v.(type) {
	case time.Time:
		return val, nil
	case *time.Time:
		if val == nil {
			return time.Time{}, uerror.New("time pointer is nil")
		}
		return *val, nil
	case int64:
		// Unix时间戳(秒)
		return time.Unix(val, 0), nil
	case int:
		// Unix时间戳(秒)
		return time.Unix(int64(val), 0), nil
	case string:
		if val == "" {
			return time.Time{}, nil
		}
		// 尝试多种格式
		formats := []string{
			time.RFC3339,          // "2006-01-02T15:04:05Z07:00"
			time.RFC3339Nano,      // "2006-01-02T15:04:05.999999999Z07:00"
			"2006-01-02 15:04:05", // DateTime
			"2006-01-02T15:04:05", // DateTime with T
			"2006-01-02",          // Date
			"2006/01/02 15:04:05", // DateTime with /
			"2006/01/02",          // Date with /
			time.RFC1123,          // "Mon, 02 Jan 2006 15:04:05 MST"
			time.RFC1123Z,         // "Mon, 02 Jan 2006 15:04:05 -0700"
		}
		for _, format := range formats {
			if t, err := time.Parse(format, val); err == nil {
				return t, nil
			}
		}
		return time.Time{}, uerror.New("unable to parse time string: " + val)
	default:
		// 尝试转换为字符串再解析
		s := ToString(v)
		return ToTime(s)
	}
}

// ToTimeDef 如果转换失败返回默认值
func ToTimeDef(v any, def time.Time) time.Time {
	if t, err := ToTime(v); err == nil {
		return t
	}
	return def
}
