package main

import (
	"fmt"

	"github.com/goburrow/modbus"
)

func main() {
	plan, err := ioutil.ReadFile(filename)
	if err != nil{
		log.Error(err)
	}
	handler := modbus.NewTCPClientHandler("localhost:502")
	// Connect manually so that multiple requests are handled in one session
	err := handler.Connect()
	defer handler.Close()
	client := modbus.NewClient(handler)

	_, err = client.WriteMultipleRegisters(0, 3, []byte{1, 3, 0, 4, 0, 5})
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	//ui := uint16(0xFF00)
	//_, err = client.WriteMultipleCoils(0, 3, []byte{1, 1, 1, 1})
	results, err := client.ReadHoldingRegisters(0, 3)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("results %v\n", results)
}
