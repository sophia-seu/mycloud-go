package main

import "fmt"

type Data struct {
	id   int
	name string
}

func getData(q interface{}) (data Data) {
	data.id = 1
	data.name = "sophia"

	return data
}

func test() (a int) {
	a = 1
	return a
}
func main() {

	res := test()
	fmt.Println(res)
}
