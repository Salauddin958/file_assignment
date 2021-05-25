package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
)

const dir = "test"

func main() {
	extensionsMap := make(map[string]int)
	ch := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go recurseDirectory(&wg, ch, dir)
	go func() {
		wg.Done()
		wg.Wait()
		close(ch)
	}()
	for v := range ch {
		extensionsMap[v]++
	}
	for key, value := range extensionsMap {
		fmt.Println(key, value)
	}
}

func recurseDirectory(wg *sync.WaitGroup, ch chan<- string, path string) {
	defer wg.Done()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			wg.Add(1)
			newPath := path + "/" + file.Name()
			go recurseDirectory(wg, ch, newPath)
		} else {
			extension := filepath.Ext(file.Name())
			ch <- extension
		}
	}

}
