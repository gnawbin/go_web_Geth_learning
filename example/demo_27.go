package main

import "fmt"

func main() {
	fmt.Println("main start")

	go func() {
		fmt.Println("goroutine")
	}()

	fmt.Println("main end")
}
