package main

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// DateOnly 是对 time.Time 的包装，专门用于映射数据库的 DATE 类型字段
// （如 invoices.invoice_date）。
//
// 解决的问题：MariaDB/MySQL 驱动在 Scan 一个 DATE 列到 *string 时，
// 会把值格式化成 RFC3339（如 "2026-05-29T00:00:00Z"）再转字符串，
// 而不是保留 "2026-05-29"。如果这个字符串未经处理直接拿去 INSERT/UPDATE，
// 会被 MySQL 以 "Incorrect date value" 错误拒绝。
//
// DateOnly 通过实现 Scanner / Valuer 接口，在数据库读写层面统一处理日期格式；
// 通过实现 MarshalJSON / UnmarshalJSON，保证序列化给前端、从前端反序列化时
// 永远是干净的 "2026-05-29" 格式，前端完全无感知，不需要做任何改动。
//
// 使用方式：把结构体里日期字段的类型从 *string 改为 *DateOnly 即可，
// 不再需要在任何 handler 里手动调用日期格式化函数。
type DateOnly struct {
	time.Time
}

// dateonly.go 顶部加这几行，相当于"显式声明"
var _ driver.Valuer = DateOnly{}     // 我实现了 Valuer
var _ sql.Scanner = &DateOnly{}      // 我实现了 Scanner
var _ json.Marshaler = DateOnly{}    // 我实现了 Marshaler
var _ json.Unmarshaler = &DateOnly{} // 我实现了 Unmarshaler

const dateLayout = "2006-01-02"

// Scan 实现 database/sql.Scanner 接口，从数据库读取 DATE 列时自动调用
func (d *DateOnly) Scan(value interface{}) error {
	if value == nil {
		d.Time = time.Time{}
		return nil
	}
	switch v := value.(type) {
	case time.Time:
		d.Time = v
		return nil
	case []byte:
		t, err := time.Parse(dateLayout, string(v))
		if err != nil {
			return fmt.Errorf("DateOnly.Scan: 解析失败: %w", err)
		}
		d.Time = t
		return nil
	case string:
		t, err := time.Parse(dateLayout, v)
		if err != nil {
			return fmt.Errorf("DateOnly.Scan: 解析失败: %w", err)
		}
		d.Time = t
		return nil
	}
	return fmt.Errorf("DateOnly.Scan: 不支持的类型 %T", value)
}

// Value 实现 database/sql/driver.Valuer 接口，写入数据库 DATE 列时自动调用
func (d DateOnly) Value() (driver.Value, error) {
	if d.Time.IsZero() {
		return nil, nil
	}
	return d.Time.Format(dateLayout), nil
}

// MarshalJSON 实现 json.Marshaler 接口，序列化给前端时自动调用
// 输出格式固定为 "2026-05-29"，零值输出 null
func (d DateOnly) MarshalJSON() ([]byte, error) {
	if d.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + d.Time.Format(dateLayout) + `"`), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口，从前端 JSON 解析时自动调用
// 兼容 "2026-05-29" 和 "2026-05-29T00:00:00Z" 两种输入格式，
// 后者用于兼容数据库 Scan 后未经处理直接传回前端再传回来的边界情况
func (d *DateOnly) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "null" || s == "" {
		d.Time = time.Time{}
		return nil
	}
	if len(s) >= 10 {
		s = s[:10] // 兼容带时间戳的输入，只取日期部分
	}
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return fmt.Errorf("DateOnly.UnmarshalJSON: 解析失败: %w", err)
	}
	d.Time = t
	return nil
}

// String 方便 fmt.Println / logger 直接打印
func (d DateOnly) String() string {
	if d.Time.IsZero() {
		return ""
	}
	return d.Time.Format(dateLayout)
}
