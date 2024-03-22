package pkg

import (
	"encoding/gob"
	"fmt"
	"os"
)

func createFile(fileName string) *os.File {
	file, err := os.Create("data/" + fileName + ".bin")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil
	}

	return file
}

func removeFile(fileName string) {
	err := os.Remove("data/" + fileName + ".bin")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
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

func Reset(filename string) {
	removeFile(filename)
}
