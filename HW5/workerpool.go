package main

import (
	"fmt"
	
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for job := range jobs {
		fmt.Println("Worker ", id, " worked ", job)
		fmt.Println("Worker ", id, " has completed job ", job)
		results <- id
	}
}
func main() {
	jobs := make(chan int, 15)
	results := make(chan int, 15)

	for i := 1; i <= 15; i++ {
		go worker(i, jobs, results)
	}

	for j := 1; j <= 15; j++ {
		jobs <- j
	}
	close(jobs)

	for z := 1; z <= 15; z++ {
	    <-results
	}
}