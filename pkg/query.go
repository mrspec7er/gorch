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
