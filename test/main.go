package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

var wg sync.WaitGroup

func main() {
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go req()
	}
	wg.Wait()
}

func req() {
	res, err := http.Get("http://localhost:8000/hello")
	if err != nil {
		fmt.Println("get data err:", err)
		return
	}
	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("get data err:", err)
		return
	}
	fmt.Println(string(data))
	wg.Done()
}
