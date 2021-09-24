package main

import "golang.org/x/tour/reader"

type MyReader struct{}

func (m MyReader) Read(l []byte) (int, error){
	i:=0
	for i < len(l){
		l[i] = 'A'
		i++
	}
	return len(l), nil
	
}


func main() {
	reader.Validate(MyReader{})
}