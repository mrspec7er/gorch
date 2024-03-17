package main

import (
	"fmt"
	"time"

	"github.com/mrspec7er/gorch/pkg"
)

type Dummy struct {
	ID    uint
	Label string
	Value int
}

func main() {
	data := []Dummy{
		{ID: 1, Label: "Name", Value: 111},
		{ID: 2, Label: "Specialty", Value: 222},
		{ID: 3, Label: "Net Worth", Value: 333},
	}

	filename := "test_" + time.Now().GoString()
	pkg.Insert(filename, data)

	time.Sleep(2 * time.Second)

	result := pkg.Find[Dummy](filename, "Value", "222")

	fmt.Println("Result: ", result)
}
