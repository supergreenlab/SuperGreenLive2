package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func main() {
	fmt.Println("pitest ok")
	fmt.Printf("  goos=%s goarch=%s\n", runtime.GOOS, runtime.GOARCH)

	hostname, err := os.Hostname()
	if err != nil {
		fmt.Printf("  hostname=error: %v\n", err)
	} else {
		fmt.Printf("  hostname=%s\n", hostname)
	}

	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			if strings.HasPrefix(line, "model name") || strings.HasPrefix(line, "Model") || strings.HasPrefix(line, "CPU architecture") {
				fmt.Printf("  %s\n", strings.TrimSpace(line))
			}
		}
	}
}
