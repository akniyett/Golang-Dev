package main

import (
	"sync"
	"fmt"
)

func merge(cs ...<-chan int) <-chan int {
	res := make(chan int)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, i := range cs {
		go func(i <-chan int) {
			for j := range i {
				res <- j
			}
			wg.Done()
		}(i)
	}
	go func() {
		wg.Wait()
		close(res)
	}()
	return res
}

func convert(ans ... int) <-chan int {

	do := make(chan int)
	go func() {
		for _, i := range ans {
			do <- i
		}
		close(do)
	}()
	return do
}


func main() {
	ch1 := convert(0, 1, 2, 3, 4)
	ch2 := convert(5, 6, 7)
	ch3 := convert(8, 9, 10, 11)
	for i := range merge(ch1, ch2, ch3) {
		fmt.Println(i)
	}
}