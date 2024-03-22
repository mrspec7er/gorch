package main

import (
	"fmt"
	"strconv"
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
		{ID: 4, Label: "Income", Value: 444},
		{ID: 5, Label: "Status", Value: 555},
	}

	filename := "test_" + "time.Now().GoString()s"
	pkg.Insert(filename, data)

	time.Sleep(2 * time.Second)

	result := pkg.Find[Dummy](filename, "Value", 222)

	for i, v := range result {
		fmt.Println("RESULT: "+strconv.Itoa(i), v)
	}

	pkg.Reset(filename)
}
