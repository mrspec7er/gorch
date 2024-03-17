package pkg

import (
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
)

type Any interface{}

func Store(fileName string) *os.File {
	file, err := os.Create("data/" + fileName + ".bin")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil
	}

	return file
}

func Read(fileName string) *os.File {
	file, err := os.Open("data/" + fileName + ".bin")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}

	return file
}

func Insert[S any](filename string, data S) {
	file := Store(filename)
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err := encoder.Encode(data)
	if err != nil {
		fmt.Println("Error encoding map to binary:", err)
		return
	}
}

func Find[S any](filename string, key string, value any) any {
	file := Read(filename)
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var data []S
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding binary data:", err)
		return nil
	}

	for _, d := range data {
		val := reflect.ValueOf(d)
		if val.Kind() != reflect.Struct {
			fmt.Println("Error: input data is not a struct")
			return nil
		}
		fieldType := val.FieldByName(key).Type()
		if fieldType != reflect.TypeOf(value) {
			fmt.Println("Field type doesn't match value type")
			continue
		}
		fieldValue := val.FieldByName(key).Interface()
		if fieldValue == value {
			return d
		}
	}

	return nil
}
