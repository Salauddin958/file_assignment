package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func getDirectory() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter Directory path")
	path, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	path = strings.TrimSuffix(path, "\r\n")
	return path
}

func main() {
	extensionsMap := make(map[string]int)
	ch := make(chan string)
	wg := sync.WaitGroup{}
	wg.Add(2)
	dir := getDirectory()
	startTime := time.Now()
	go recursiveDirWalkthrough(&wg, ch, dir)
	go func() {
		wg.Done()
		wg.Wait()
		close(ch)
	}()
	for v := range ch {
		extensionsMap[v]++
	}
	countOfFiles := 0
	for key, value := range extensionsMap {
		countOfFiles += value
		fmt.Println(key, value)
	}
	fmt.Println("number of files : ", countOfFiles)
	fmt.Println("time taken : ", time.Since(startTime))
}

func recursiveDirWalkthrough(wg *sync.WaitGroup, ch chan<- string, path string) {
	defer wg.Done()
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			wg.Add(1)
			newPath := path + "/" + file.Name()
			go recursiveDirWalkthrough(wg, ch, newPath)
		} else {
			extension := filepath.Ext(file.Name())
			ch <- extension
		}
	}

}
