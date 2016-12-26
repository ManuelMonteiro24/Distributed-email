package main

import (
	"fmt"
	"os"
)

func main() {
	createFile("ola.txt")
	writeFile("ola","ola.txt")
}
func createFile(file_name string) {
	// detect if file exists
	var _, err = os.Stat(file_name)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(file_name)
		checkError(err)
		defer file.Close()
  }
}

func writeFile(key string, file_name string) {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(file_name, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()

	// write some text to file
	_, err = file.WriteString(key)
	checkError(err)

	// save changes
	err = file.Sync()
	checkError(err)
}
