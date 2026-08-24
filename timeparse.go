package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// 中文星期
var weekdayZH = map[time.Weekday]string{
	time.Sunday:    "星期日",
	time.Monday:    "星期一",
	time.Tuesday:   "星期二",
	time.Wednesday: "星期三",
	time.Thursday:  "星期四",
	time.Friday:    "星期五",
	time.Saturday:  "星期六",
}

// 支持的时间字符串格式（本地时区解析）
var timeLayouts = []string{
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006/01/02 15:04:05",
	"2006-01-02 15:04",
	"2006-01-02",
	"2006/01/02",
}

// Result 表示一次识别结果
type Result struct {
	Kind    string    // "timestamp10" | "timestamp13" | "datetime"
	Input   string    // 清洗后的输入
	Time    time.Time // 解析出的本地时间
	TSSec   int64     // 秒级时间戳
	TSMilli int64     // 毫秒级时间戳
}

// Parse 从选中的文本中识别时间戳或时间字符串。
// 会先做基本清洗（去空白、去引号、去逗号）。
func Parse(raw string) (*Result, error) {
	s := clean(raw)
	if s == "" {
		return nil, fmt.Errorf("未选中任何文本")
	}

	// 1) 纯数字：按位数判断秒 / 毫秒时间戳
	if isAllDigits(s) {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("数字过大，无法解析为时间戳: %s", s)
		}
		switch len(s) {
		case 8:
			if t, err := time.ParseInLocation("20060102", s, time.Local); err == nil {
				return &Result{
					Kind:    "datetime",
					Input:   s,
					Time:    t,
					TSSec:   t.Unix(),
					TSMilli: t.UnixMilli(),
				}, nil
			}
		case 10:
			t := time.Unix(n, 0)
			return &Result{Kind: "timestamp10", Input: s, Time: t, TSSec: n, TSMilli: n * 1000}, nil
		case 13:
			t := time.UnixMilli(n)
			return &Result{Kind: "timestamp13", Input: s, Time: t, TSSec: n / 1000, TSMilli: n}, nil
		default:
			return nil, fmt.Errorf("数字位数为 %d，仅支持 10 位（秒）或 13 位（毫秒）时间戳: %s", len(s), s)
		}
	}

	// 2) 时间字符串：按候选布局解析（本地时区）
	for _, layout := range timeLayouts {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return &Result{
				Kind:    "datetime",
				Input:   s,
				Time:    t,
				TSSec:   t.Unix(),
				TSMilli: t.UnixMilli(),
			}, nil
		}
	}

	return nil, fmt.Errorf("无法识别的时间格式: %q\n支持：10/13 位时间戳，或 2006-01-02 15:04:05", s)
}

// Format 生成人类可读的多行结果文本
func (r *Result) Format() string {
	t := r.Time
	wd := weekdayZH[t.Weekday()]
	var b strings.Builder
	fmt.Fprintf(&b, "识别输入：%s\n", r.Input)
	fmt.Fprintf(&b, "类型：%s\n", kindZH(r.Kind))
	b.WriteString("──────────────\n")
	fmt.Fprintf(&b, "时间：%s\n", t.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&b, "星期：%s\n", wd)
	fmt.Fprintf(&b, "秒级时间戳(10位)：%d\n", r.TSSec)
	fmt.Fprintf(&b, "毫秒时间戳(13位)：%d\n", r.TSMilli)
	fmt.Fprintf(&b, "时区：%s", t.Format("-07:00 MST"))
	return b.String()
}

// FormatShort 用于通知气泡的紧凑单行/多行文本
func (r *Result) FormatShort() string {
	t := r.Time
	return fmt.Sprintf("%s %s\n秒:%d 毫秒:%d",
		t.Format("2006-01-02 15:04:05"), weekdayZH[t.Weekday()], r.TSSec, r.TSMilli)
}

func kindZH(k string) string {
	switch k {
	case "timestamp10":
		return "10 位秒级时间戳"
	case "timestamp13":
		return "13 位毫秒时间戳"
	case "datetime":
		return "时间字符串"
	}
	return k
}

func clean(raw string) string {
	s := strings.TrimSpace(raw)
	// 去掉常见包裹字符
	s = strings.Trim(s, "\"'`,;")
	s = strings.TrimSpace(s)
	return s
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
