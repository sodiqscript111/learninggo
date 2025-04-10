package main

import (
	"fmt"
	"time"
)

func countTo(n int) {
	for i := 1; i < n+1; i++ {
		fmt.Printf("%d\n", i)
	}
}

func main() {
	fmt.Println("Main function is running")
	go countTo(5)
	time.Sleep(3 * time.Second)
}
