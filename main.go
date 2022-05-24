package main

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/go-xmlfmt/xmlfmt"
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
		go processJson(data, ch, &wg)
		go processXml(data, ch, &wg)
		wg.Wait()
		elapsedTime := time.Since(start)
		log.Printf("execution time %s\n", elapsedTime)
	}
}

func processJson(data []byte, ch <-chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	buffer, err := formatJson(data)
	if err != nil {
		log.Default().Println("Error while checking JSON: ", err)
		return
	}
	clipboard.Write(clipboard.FmtText, buffer.Bytes())
	<-ch // we want to ignore the message that has just been written, as it's the formatted code
}

func formatJson(jsonBytes []byte) (bytes.Buffer, error) {
	var out bytes.Buffer
	if !json.Valid(jsonBytes) {
		return out, fmt.Errorf("is not JSON")
	}
	err := json.Indent(&out, jsonBytes, "", identTwoSpaces)
	return out, err
}

func processXml(data []byte, ch <-chan []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	formattedXml, err := formatXml(data)
	if err != nil {
		log.Default().Println("Error while checking XML: ", err)
		return
	}
	clipboard.Write(clipboard.FmtText, []byte(formattedXml))
	<-ch // we want to ignore the message that has just been written, as it's the formatted code
}

func formatXml(data []byte) (string, error) {
	if !isValidXml(data) {
		return "", fmt.Errorf("is not XML")
	}
	formattedXml := xmlfmt.FormatXML(string(data), "", identTwoSpaces, true)
	return formattedXml, nil
}

func isValidXml(data []byte) bool {
	// is not necessary to unmarshal it, however it's the easiest way to check if it's valid
	return xml.Unmarshal(data, new(any)) == nil
}
