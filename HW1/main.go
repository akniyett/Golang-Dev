package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z:=1.0
	i := 1
	for {
		if math.Abs(z - z - (z * z - x) / (2 * z)) <= 1e-9 {
			break
		}
		z -= (z*z - x) / (2*z)
		fmt.Println(i, z)
		i++
	}
	return z
	
}

func main() {
	fmt.Println(Sqrt(17))
	fmt.Println(math.Sqrt(17))
}
