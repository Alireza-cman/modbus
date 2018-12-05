package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/tbrandon/mbserver"
)

func main() {
	serv := mbserver.NewServer()
	err := serv.ListenTCP("127.0.0.1:1502")
	log.Println("[Zooring] trying to listen to the port")
	if err != nil {
		log.Printf("%v\n", err)
	}

	log.Println("[done] server is connected ")

	defer serv.Close()

	// Wait forever
	for {
		time.Sleep(5 * time.Second)
		fmt.Println(serv.HoldingRegisters[0:10])
	}

	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()

}
