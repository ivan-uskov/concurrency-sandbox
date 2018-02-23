package main

import (
	"fmt"
)

func numbers(done chan struct{}) <-chan int {
	ch := make(chan int)
	go func() {
		i := 0
		for {
			select {
				case <- done:
					close(ch)
					return
				case ch <- func() int { i++; return i }():
			}
		}
	}()
	return ch
}

func main() {
	doneCh := make(chan struct{})
	for i := range numbers(doneCh) {
		fmt.Printf("New number %d\n", i)
		if i == 7 {
			doneCh <- struct{}{}
		}
	}
}