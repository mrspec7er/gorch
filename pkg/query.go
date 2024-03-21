package pkg

import (
	"encoding/gob"
	"fmt"
	"os"
	"reflect"
	"sync"
)

func readFile(fileName string) (*os.File, error) {
	file, err := os.OpenFile("data/"+fileName+".bin", 0666, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func filterData[S any](data S, key string, value any, ctx chan []*S, wg *sync.WaitGroup) {
	field := reflect.ValueOf(data)
	if field.Kind() != reflect.Struct {
		fmt.Println("Error: input data is not a struct")
		wg.Done()
		return
	}
	fieldType := field.FieldByName(key).Type()
	if fieldType != reflect.TypeOf(value) {
		fmt.Println("Field type doesn't match value type")
		wg.Done()
		return
	}
	storedValue := field.FieldByName(key).Interface()
	if storedValue == value {
		ctx <- []*S{&data}
		wg.Done()
		return
	}

	wg.Done()
}

func Find[S any](filename string, key string, value any) []*S {
	file, err := readFile(filename)
	if err != nil {
		fmt.Println("Error decoding binary data:", err)
		return nil
	}

	defer file.Close()

	decoder := gob.NewDecoder(file)
	var storedData []S
	err = decoder.Decode(&storedData)
	if err != nil {
		fmt.Println("Error decoding binary storedData:", err)
		return nil
	}

	var finalResult []*S

	ctx := make(chan []*S, len(storedData))
	wg := &sync.WaitGroup{}

	wg.Add(len(storedData))
	for _, d := range storedData {
		go filterData(d, key, value, ctx, wg)
	}

	wg.Wait()
	close(ctx)

	for v := range ctx {
		finalResult = append(finalResult, v...)
	}

	return finalResult
}
