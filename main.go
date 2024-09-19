package main

import (
	"flag"
	"fmt"
	"ping/ping"
)

func main() {
	// 定义命令行参数，支持缩写
	count := flag.Int("c", 4, "要发送的回显请求数量")
	interval := flag.Int("i", 1, "发送每个数据包之间的等待间隔（秒）")
	size := flag.Int("s", 32, "要发送的数据包大小（字节）")
	ttl := flag.Int("t", 64, "数据包的生存时间（TTL）")
	timeout := flag.Int("w", 2, "等待每个回复的超时时间（秒）")

	// 定义全写参数
	flag.IntVar(count, "count", *count, "要发送的回显请求数量")
	flag.IntVar(interval, "interval", *interval, "发送每个数据包之间的等待间隔（秒）")
	flag.IntVar(size, "size", *size, "要发送的数据包大小（字节）")
	flag.IntVar(ttl, "ttl", *ttl, "数据包的生存时间（TTL）")
	flag.IntVar(timeout, "timeout", *timeout, "等待每个回复的超时时间（秒）")

	flag.Parse()

	// 检查参数并获取目标地址
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: go run main.go [options] <address>")
		flag.PrintDefaults()
		return
	}
	address := args[len(args)-1] // 获取最后一个参数作为地址

	// 打印解析的参数
	fmt.Printf("参数: count=%d, interval=%d, size=%d, ttl=%d, timeout=%d, address=%s\n",
		*count, *interval, *size, *ttl, *timeout, address)

	// 调用 SendPing 函数
	ping.SendPing(address, *count, *interval, *size, *timeout, uint8(*ttl))
}
