package main

import (
	"fmt"
	"math"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string{
	return fmt.Sprint("cannot Sqrt negative number: ", float64(e))
	
}

func Sqrt(x float64) (float64, error) {
	if x < 0{
		return 0, ErrNegativeSqrt(x)
	}else if x == 0{
		return 0, nil
	}
	
	z:=1.0
	for {
		if math.Abs(z - z - (z * z - x) / (2 * z)) <= 1e-9 {
			break
		}
		z -= (z*z - x) / (2*z)
	}
	return z, nil



}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
	fmt.Println(Sqrt(0))
}