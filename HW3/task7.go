package main

import (
	"golang.org/x/tour/tree"
	"fmt"
)


func Walk(t *tree.Tree, ch chan int) {
	if t != nil {
		if t.Left != nil {
			Walk(t.Left, ch)
		}
		ch<-t.Value
		if t.Right != nil {
			Walk(t.Right, ch)
		}
	}
	
	
}

func Same(t1, t2 *tree.Tree) bool {
 	ch1 := make(chan int)
	ch2 := make(chan int)
 	go Walk(t1, ch1)
 	go Walk(t2, ch2)
	i := 0
	for i < 8 {
		if <-ch1 != <-ch2  {
   			return false
  		}
		i++
 	}
	return true
}

func main() {
	
	ch := make(chan int, 8)
	go Walk(tree.New(1), ch)
	
	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}