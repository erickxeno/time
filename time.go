// @Package: time
// @Filename: time.go
// @Author: xeno
// @Create date: 2025.06.06
// @Description: Shared time.Now() function, avoiding every system call

package time

import (
	"strings"
	"sync/atomic"
	osTime "time"
	"unsafe"
)

const (
	zeroAscii = '0'
)

var (
	currentTime *Time
	clock       = osTime.Millisecond * 1 // init default clock interval as 1ms
	ticker      = osTime.NewTicker(clock)
	zone        []byte
)

type Time struct {
	osTime.Time
	serialBytes []byte // the format string for now time
}

// SetClock set the refresh interval to new duration
func SetClock(t osTime.Duration) {
	atomic.StoreInt64((*int64)(&clock), int64(t))
}

type TimePrecision int

const (
	TimePrecisionSecond TimePrecision = iota
	TimePrecisionMillisecond
	TimePrecisionMicrosecond
)

var (
	timePrecision = TimePrecisionMillisecond // default time precision is millisecond
)

func SetTimePrecision(tm TimePrecision) {
	if tm == TimePrecisionMicrosecond {
		localC := atomic.LoadInt64((*int64)(&clock))
		if localC > int64(osTime.Microsecond) {
			SetClock(osTime.Microsecond * 1)
		}
	}
	atomic.StoreInt32((*int32)(unsafe.Pointer(&timePrecision)), int32(tm))
}

func Current() Time {
	return *(*Time)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&currentTime))))
}

func Now() osTime.Time {
	return Current().Time
}

func (t *Time) String() string {
	return string(t.serialBytes[:len(t.serialBytes)-5])
}

func (t *Time) StringWithZone() string {
	return string(t.serialBytes)
}

func (t *Time) ReadOnlyDataWithoutZone() []byte {
	return t.serialBytes[:len(t.serialBytes)-5]
}

func (t *Time) ReadOnlyDataWithZone() []byte {
	return t.serialBytes
}

func refreshTask() {
	localC := atomic.LoadInt64((*int64)(&clock))
	for {
		cur := <-ticker.C
		refreshCurrentTime(cur)
		clock := atomic.LoadInt64((*int64)(&clock))
		if clock != localC {
			ticker.Stop()
			ticker = osTime.NewTicker(osTime.Duration(clock))
		}
	}
}

func refreshCurrentTime(cur osTime.Time) {
	curT := Time{
		Time: cur,
	}
	precision := atomic.LoadInt32((*int32)(unsafe.Pointer(&timePrecision)))
	switch precision {
	case int32(TimePrecisionSecond):
		curT.serialBytes = make([]byte, 0, 24)
	case int32(TimePrecisionMillisecond):
		curT.serialBytes = make([]byte, 0, 28)
	case int32(TimePrecisionMicrosecond):
		curT.serialBytes = make([]byte, 0, 31)
	}
	curT.serialBytes = timeData(cur, curT.serialBytes, TimePrecision(precision))
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&currentTime)), unsafe.Pointer(&curT))
}

func init() {
	timeStr := strings.Split(osTime.Now().String(), " ") // 2023-08-12 16:07:33.282171 +0800 CST m=+0.391647800
	if len(timeStr) >= 3 {
		zone = []byte(timeStr[2])
	}
	refreshCurrentTime(osTime.Now())
	go refreshTask()
}

func timeData(t osTime.Time, c []byte, precision TimePrecision) []byte {
	year, month, day := t.Date()
	// year
	c = append(c, byte(year/1000)+zeroAscii)
	c = append(c, byte(year%1000/100)+zeroAscii)
	c = append(c, byte(year%100/10)+zeroAscii)
	c = append(c, byte(year%10)+zeroAscii)
	c = append(c, '-')
	// month
	c = append(c, byte(month/10)+zeroAscii)
	c = append(c, byte(month%10)+zeroAscii)
	c = append(c, '-')
	// day
	c = append(c, byte(day/10)+zeroAscii)
	c = append(c, byte(day%10)+zeroAscii)
	c = append(c, ' ')

	hour, min, sec := t.Clock()
	// hour
	c = append(c, byte(hour/10)+zeroAscii)
	c = append(c, byte(hour%10)+zeroAscii)
	c = append(c, ':')
	// min
	c = append(c, byte(min/10)+zeroAscii)
	c = append(c, byte(min%10)+zeroAscii)
	c = append(c, ':')
	// min
	c = append(c, byte(sec/10)+zeroAscii)
	c = append(c, byte(sec%10)+zeroAscii)

	switch precision {
	case TimePrecisionSecond:
	case TimePrecisionMillisecond:
		c = append(c, ',')
		ms := t.Nanosecond() / 1e6
		c = append(c, byte(ms/100)+zeroAscii)
		c = append(c, byte(ms%100/10)+zeroAscii)
		c = append(c, byte(ms%10)+zeroAscii)
	case TimePrecisionMicrosecond:
		c = append(c, ',')
		us := t.Nanosecond() / 1e3
		c = append(c, byte(us/100000)+zeroAscii)
		c = append(c, byte(us%100000/10000)+zeroAscii)
		c = append(c, byte(us%10000/1000)+zeroAscii)
		c = append(c, byte(us%1000/100)+zeroAscii)
		c = append(c, byte(us%100/10)+zeroAscii)
		c = append(c, byte(us%10)+zeroAscii)
	}

	// zone
	if len(zone) != 0 {
		c = append(c, zone...)
	}
	return c
}
