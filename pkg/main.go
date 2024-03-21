package pkg

import (
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
	"sync"
)

func createFile(fileName string) *os.File {
	file, err := os.Create("data/" + fileName + ".bin")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil
	}

	return file
}

func readFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile("data/"+fileName+".bin", 0666, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func Insert[S any](filename string, data []S) {
	file, err := readFile(filename)
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
		file := createFile(filename)
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

func Find[S any](filename string, key string, value any) []*S {
	file, err := readFile(filename)
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

	var finalResult []*S

	ctx := make(chan []*S, len(data))
	wg := &sync.WaitGroup{}

	wg.Add(len(data))
	for _, d := range data {
		go filterData(d, key, value, ctx, wg)
	}

	wg.Wait()
	close(ctx)

	for v := range ctx {
		finalResult = append(finalResult, v...)
	}

	return finalResult
}

func filterData[S any](d S, key string, value any, ctx chan []*S, wg *sync.WaitGroup) {
	val := reflect.ValueOf(d)
	if val.Kind() != reflect.Struct {
		fmt.Println("Error: input data is not a struct")
		wg.Done()
		return
	}
	fieldType := val.FieldByName(key).Type()
	if fieldType != reflect.TypeOf(value) {
		fmt.Println("Field type doesn't match value type")
		wg.Done()
		return
	}
	fieldValue := val.FieldByName(key).Interface()
	if fieldValue == value {
		ctx <- []*S{&d}
		wg.Done()
		return
	}

	wg.Done()
}
