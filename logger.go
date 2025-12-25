package main

import (
	"fmt"
	"sync"
)

type logger struct {
	loggerChan chan string
	wg         sync.WaitGroup
}

func (obj *logger) Init() {
	obj.loggerChan = make(chan string, 100) // this is unbuffered as of now
	go func() {
		for message := range obj.loggerChan {
			obj.PrintMessage(message)
		}
	}()
}

func (obj *logger) PrintMessage(message string) {
	defer obj.wg.Done()
	fmt.Println(message)
}

func (obj *logger) SendLog(message string) {
	obj.wg.Add(1)
	obj.loggerChan <- message
}

func (obj *logger) GraceFulShutdown() {
	close(obj.loggerChan)
	obj.wg.Wait()
}