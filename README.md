# Time

一个高性能的Go语言时间处理库，通过缓存系统时间并提供高效的时间格式化功能，显著减少系统调用开销。

## 特性

- 🚀 高性能：通过缓存系统时间，避免频繁的系统调用
- ⚡ 可配置：支持自定义时间刷新间隔
- 🔒 线程安全：使用原子操作确保并发安全
- 📝 多种格式：提供多种时间格式化输出方法
- 🎯 高效序列化：优化的时间序列化处理

## 性能测试

在Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz上的基准测试结果：

```
BenchmarkTimer/ticker(default_1ms_ticker)-12            1000000000               0.1537 ns/op       0 B/op           0 allocs/op
BenchmarkTimer/std-12                                   140642682                8.712 ns/op        0 B/op           0 allocs/op
BenchmarkTimer/ticker(10ms_ticker)-12                   1000000000               0.1634 ns/op       0 B/op           0 allocs/op
BenchmarkTimer/ticker(1us_ticker)-12                    1000000000               0.1655 ns/op       0 B/op           0 allocs/op
```

性能对比：
- 标准库时间获取：约 8.71 ns/op
- 本库时间获取（不同ticker间隔）：
  - 默认1ms间隔：约 0.15 ns/op
  - 10ms间隔：约 0.16 ns/op
  - 1μs间隔：约 0.17 ns/op
- 性能提升：约 57 倍（相比标准库）

测试环境：
- Go版本：1.21.7
- 操作系统：macOS (darwin)
- CPU：Intel(R) Core(TM) i7-9750H @ 2.60GHz

## 安装

```bash
go get github.com/erickxeno/time
```

## 使用示例

### 基本使用

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    // 获取当前时间
    now := time.Now()
    fmt.Println(now)

    // 获取当前时间（自定义类型）
    current := time.Current()
    fmt.Println(current)
}
```

### 自定义刷新间隔

```go
package main

import (
    "time"
)

func main() {
    // 设置时间刷新间隔为100毫秒
    time.SetClock(time.Millisecond * 100)
}
```

### 时间格式化输出

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    current := time.Current()
    
    // 获取不带时区的时间字符串
    timeStr := current.String()
    fmt.Println(timeStr)  // 输出格式：2024-03-21 15:04:05,123
    
    // 获取带时区的时间字符串
    timeStrWithZone := current.StringWithZone()
    fmt.Println(timeStrWithZone)  // 输出格式：2024-03-21 15:04:05,123 +0800 CST
    
    // 获取时间字节数组（不带时区）
    timeBytes := current.ReadOnlyDataWithoutZone()
    fmt.Println(string(timeBytes))
    
    // 获取时间字节数组（带时区）
    timeBytesWithZone := current.ReadOnlyDataWithZone()
    fmt.Println(string(timeBytesWithZone))
}
```

## 性能优势

该库通过以下方式优化性能：

1. 缓存系统时间，减少系统调用
2. 使用原子操作确保并发安全
3. 预分配内存，减少内存分配
4. 优化的时间格式化算法

## 注意事项

- 默认时间刷新间隔为1毫秒
- 时间格式化输出固定为 "YYYY-MM-DD HH:mm:ss,SSS" 格式
- 时区信息在程序启动时确定，运行期间不会改变

## License

MIT License