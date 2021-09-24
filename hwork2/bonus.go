package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (rt rot13Reader) Read (list []byte) (int, error){
	a, err := rt.r.Read(list)
	i := 0
	for i < len(list){
		if (list[i] >= 'a' && list[i] <= 'm') || (list[i] >= 'A' && list[i] <= 'M'){
			list[i] = list[i] + 13
	}else if (list[i] <= 'z' && list[i] >= 'n') || (list[i] <= 'Z' && list[i] >= 'N'){
		list[i] =list[i] - 13
	}
	
	i++	
	}
	return a, err
	
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
