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

func Read(fileName string) (*os.File, error) {
	file, err := os.OpenFile("data/"+fileName+".bin", 0666, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func Insert[S any](filename string, data []S) {
	file, err := Read(filename)
	if err == nil {
		decoder := gob.NewDecoder(file)
		var currentData []S
		err := decoder.Decode(&currentData)
		if err != nil {
			fmt.Println("Error decoding binary data:", err)
			return
		}

		currentData = append(currentData, data...)
		fmt.Println(currentData)

		_, err = file.Seek(0, 0)
		if err != nil {
			fmt.Println("Error seeking file:", err)
			return
		}

		encoder := gob.NewEncoder(file)
		err = encoder.Encode(currentData)
		if err != nil {
			fmt.Println("Error encoding map to binary:", err)
			return
		}

	} else {
		file := Store(filename)
		defer file.Close()

		encoder := gob.NewEncoder(file)
		err := encoder.Encode(data)
		if err != nil {
			fmt.Println("Error encoding map to binary:", err)
			return
		}

	}
	defer file.Close()
}

func Find[S any](filename string, key string, value any) *S {
	file, err := Read(filename)
	if err != nil {
		fmt.Println("Error decoding binary data:", err)
		return nil
	}

	defer file.Close()

	decoder := gob.NewDecoder(file)
	var data []S
	err = decoder.Decode(&data)
	if err != nil {
		fmt.Println("Error decoding binary data:", err)
		return nil
	}

	var finalResult *S

	for _, d := range data {
		filterResult := FilterData(d, key, value)
		if filterResult != nil {
			finalResult = filterResult
		}
	}

	return finalResult
}

func FilterData[S any](d S, key string, value any) *S {
	val := reflect.ValueOf(d)
	if val.Kind() != reflect.Struct {
		fmt.Println("Error: input data is not a struct")
		return nil
	}
	fieldType := val.FieldByName(key).Type()
	if fieldType != reflect.TypeOf(value) {
		fmt.Println("Field type doesn't match value type")
		return nil
	}
	fieldValue := val.FieldByName(key).Interface()
	if fieldValue == value {
		return &d
	}

	return nil
}
