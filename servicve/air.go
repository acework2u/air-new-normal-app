package main

import (
	"log"
	"sync"
	"time"
)

var wg sync.WaitGroup

//func main() {
//	fmt.Println("Start Main")
//	wg.Add(2)
//	go side()
//	go side()
//	fmt.Println("Return to main")
//	wg.Wait()
//	time.Sleep(5 * time.Second)
//	fmt.Println("End Main")
//}
//
//func side() {
//	fmt.Println("Start side process")
//	time.Sleep(1 * time.Second)
//	fmt.Println("End side process")
//	wg.Done()
//}

func main() {
	log.Println("Start Main")
	wg.Add(2)
	//side()
	go side()
	go side()
	//go side()
	log.Println("Return to main")
	wg.Wait()
	time.Sleep(2 * time.Second)
	log.Println("End Main")
}

func side() {
	log.Println("Start side process")
	time.Sleep(1 * time.Second)
	log.Println("End side process")
	wg.Done()
}

//var message = make(chan string)
//
//func main() {
//	go createPing()
//	msg := <-message
//	fmt.Println(msg)
//
//}
//func createPing() {
//	message <- "ping"
//}
