package main

import (
	"fmt"
	"os"
	"ping/ping"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <address> <count>")
		return
	}
	address := os.Args[1]
	count := 0
	fmt.Sscanf(os.Args[2], "%d", &count)
	ping.SendPing(address, count)
}
