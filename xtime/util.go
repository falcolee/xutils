package xtime

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrCanceled is canceled error
	ErrCanceled = errors.New("xtime: canceled")
	// ErrTimeouted is timeouted error
	ErrTimeouted = errors.New("xtime: timeouted")
)

// TimeCallback is a callback with one return value
type TimeCallback func() interface{}

// Now returns time.Now
func Now() time.Time {
	return time.Now()
}

// String returns string of now
func String() string {
	return TimeToStr(S())
}

// S returns unix timestamp in Second
func S() int64 {
	return Now().Unix()
}

// Ns returns unix timestamp in Nanosecond
func Ns() int64 {
	return Now().UnixNano()
}

// Us returns unix timestamp in Microsecond
func Us() int64 {
	return Now().UnixNano() / int64(time.Microsecond)
}

// Ms returns unix timestamp in Millisecond
func Ms() int64 {
	return Now().UnixNano() / int64(time.Millisecond)
}

// Sleep n Second
func Sleep(n int64) {
	time.Sleep(time.Duration(n) * time.Second)
}

// Usleep n Microsecond
func Usleep(n int64) {
	time.Sleep(time.Duration(n) * time.Microsecond)
}

// StrToTime returns unix timestamp of time string
func StrToTime(s string, layout ...string) (int64, error) {
	format := "2006-01-02 15:04:05"
	if len(layout) > 0 && layout[0] != "" {
		format = layout[0]
	} else {
		if len(s) == 10 {
			format = format[:10]
		}
	}

	t, err := time.ParseInLocation(format, s, time.Local)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}

// TimeToStr returns time string of unix timestamp, format in time.Local
func TimeToStr(n int64, layout ...string) string {
	format := "2006-01-02 15:04:05"
	if len(layout) > 0 && layout[0] != "" {
		format = layout[0]
	}

	return time.Unix(n, 0).Format(format)
}

// WithTimeout execute the callback with timeout return a chan and cancel func
func WithTimeout(fn TimeCallback, timeout time.Duration) (chan interface{}, func()) {
	q := make(chan bool)
	r := make(chan interface{})

	go func() {
		r <- fn()
	}()

	go func() {
		t := time.After(timeout)
		select {
		case <-t:
			r <- ErrTimeouted
			return
		case <-q:
			r <- ErrCanceled
			return
		}
	}()

	return r, func() { close(q) }
}

// SetTimeout execute the callback after timeout return a chan and cancel func
func SetTimeout(fn TimeCallback, timeout time.Duration) (chan interface{}, func()) {
	q := make(chan bool)
	r := make(chan interface{})

	go func() {
		t := time.After(timeout)
		select {
		case <-t:
			r <- fn()
		case <-q:
			r <- ErrCanceled
			return
		}
	}()

	return r, func() { close(q) }
}

// SetInterval execute the callback every timeout return a chan and cancel func
func SetInterval(fn TimeCallback, timeout time.Duration) (chan interface{}, func()) {
	q := make(chan bool)
	r := make(chan interface{})

	go func() {
		t := time.NewTicker(timeout)
		for {
			select {
			case <-t.C:
				go func() { r <- fn() }()
			case <-q:
				t.Stop()
				r <- ErrCanceled
				return
			}
		}
	}()

	return r, func() { close(q) }
}

// FormatLayout ...
func FormatLayout(layout string) string {
	layout = strings.ToLower(layout)
	for k, v := range formatMap {
		layout = strings.ReplaceAll(layout, k, v)
	}
	return layout
}

// Format ...
func Format(t time.Time, layouts ...string) string {
	layout := FormatTime
	if len(layouts) > 0 {
		layout = layouts[0]
	}
	return t.Format(FormatLayout(layout))
}

// Parse ...
func Parse(t string, layouts ...string) (time.Time, error) {
	layout := FormatTime
	if len(layouts) > 0 {
		layout = layouts[0]
	}
	location, _ := time.LoadLocation("Asia/Shanghai")
	return time.ParseInLocation(FormatLayout(layout), t, location)
}

// First ...
func First(t time.Time) int64 {
	if t.IsZero() {
		t = time.Now()
	}
	timeStr := fmt.Sprintf("%s 00:00:00", Format(t, FormatDateBar))
	tt, _ := Parse(timeStr, FormatTime)
	return tt.Unix()
}

// Last ...
func Last(t time.Time) int64 {
	return First(t) + 24*3600 - 1
}

// Check 校验日期有效性
func Check(year, month, day int) bool {
	if month < 1 || month > 12 || day < 1 || day > 31 || year < 1 || year > 32767 {
		return false
	}
	switch month {
	case 4, 6, 9, 11:
		if day > 30 {
			return false
		}
	case 2:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			if day > 29 {
				return false
			}
		} else if day > 28 {
			return false
		}
	}
	return true
}

// IsLeapYear 是否是闰年
func IsLeapYear(year int) bool {
	return (year%4 == 0 && year%100 != 0) || (year%400 == 0)
}
