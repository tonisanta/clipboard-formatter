package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/clbanning/mxj/v2"
	"golang.design/x/clipboard"
	"log"
	"sync"
	"time"
)

const identTwoSpaces = "  "

func main() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}

	fmt.Println("Listening ...")

	ch := clipboard.Watch(context.Background(), clipboard.FmtText)
	var wg sync.WaitGroup // number of working goroutines
	for data := range ch {
		log.Println("new data to process")
		wg.Add(2)
		start := time.Now()
		go processData(data, ch, &wg, formatJson)
		go processData(data, ch, &wg, formatXml)
		wg.Wait()
		elapsedTime := time.Since(start)
		log.Printf("execution time %s\n", elapsedTime)
	}
}

func formatJson(data []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", identTwoSpaces)
	return out.Bytes(), err
}

func formatXml(data []byte) ([]byte, error) {
	return mxj.BeautifyXml(data, "", identTwoSpaces)
}

func processData(data []byte, ch <-chan []byte, wg *sync.WaitGroup, format func([]byte) ([]byte, error)) {
	defer wg.Done()
	out, err := format(data)
	if err != nil {
		return
	}
	clipboard.Write(clipboard.FmtText, out)
	<-ch // we want to ignore the message that has just been written, as it's the formatted code
}
