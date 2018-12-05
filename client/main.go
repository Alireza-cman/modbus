package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Atrovan/Modbus/client/variable"
	"github.com/op/go-logging"
	"github.com/tidwall/gjson"

	"github.com/goburrow/modbus"
)

var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {

	//logger format, dont touch it
	backend1 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2 := logging.NewLogBackend(os.Stderr, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend2, format)
	backend1Leveled := logging.AddModuleLevel(backend1)
	backend1Leveled.SetLevel(logging.ERROR, "")
	logging.SetBackend(backend1Leveled, backend2Formatter)

	// reading the config json file
	configFile := variable.ConfigFile
	plan, err := ioutil.ReadFile(configFile)
	configFile = string(plan)
	if err != nil {
		log.Error(err)
		return
	}
	//retreiveing number of interfaces
	interfaceNumber := gjson.Get(configFile, "servers.#").Num
	for i := 0; i < int(interfaceNumber); i++ {

		connection := gjson.Get(configFile, "servers."+strconv.Itoa(i)+".transport.type").Str
		connection = strings.ToLower(connection)
		//fmt.Println(gjson.Get(configFile, "servers."+strconv.Itoa(i)).Raw)
		fmt.Println(connection)
		if connection == "tcp" {

			go modbusManipulationTCP(gjson.Get(configFile, "servers."+strconv.Itoa(i)).Raw)

		}

	}
	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()

}

// this function retrieve the host and port of the slave device, moreover, it retrieves the timeout which is necessary
// for the connection. at the end of the function, telemetry values which is needed for us will be fetched
func modbusManipulationTCP(configFile string) {
	host := gjson.Get(configFile, "transport.host").Str
	port := gjson.Get(configFile, "transport.port").Num
	timeout := gjson.Get(configFile, "transport.timeout").Num
	if timeout == 0 {
		timeout = 1
	}
	if port == 0 {
		port = 502
	}

	address := host + ":" + strconv.Itoa(int(port))
	handler := modbus.NewTCPClientHandler(address)
	handler.Timeout = time.Duration(timeout) * time.Second
	log.Warning("trying to make a TCP modbus connection,", address)
	//
	deviceNumber := gjson.Get(configFile, "devices.#").Num
	log.Warning("total number of devices is:", deviceNumber)
	for i := 0; i < int(deviceNumber); i++ {
		unitId := gjson.Get(configFile, "devices."+strconv.Itoa(i)+".unitId").Num
		deviceName := gjson.Get(configFile, "devices."+strconv.Itoa(i)+".deviceName").Str
		telemetry := gjson.Get(configFile, "devices."+strconv.Itoa(i)+".timeseries").Raw
		go telemetryHandler(telemetry, unitId, deviceName, handler)
		//fmt.Println(telemetry)
		//fmt.Println(unitId, deviceName)

	}

	// var wg sync.WaitGroup

	// wg.Add(1)
	// wg.Wait()

}

func telemetryHandler(json string, unitId float64, deviceName string, handler *modbus.TCPClientHandler) {
	telemetryNumber := gjson.Get(json, "#").Num
	log.Warning("telemetry number of device ", deviceName, "is:", telemetryNumber)
	for i := 0; i < int(telemetryNumber); i++ {
		tag := gjson.Get(json, strconv.Itoa(i)+".tag").Str
		kind := gjson.Get(json, strconv.Itoa(i)+".kind").Str
		functionCode := gjson.Get(json, strconv.Itoa(i)+".functionCode").Num
		address := gjson.Get(json, strconv.Itoa(i)+".address").Num
		count := gjson.Get(json, strconv.Itoa(i)+".count").Num
		period := gjson.Get(json, strconv.Itoa(i)+".pollPeriod").Num
		if period == 0 {
			period = 1
		}
		if count == 0 {
			count = 2
		}
		go atom(tag, kind, uint16(functionCode), uint16(address), uint16(count), uint16(period), handler)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()

	//functionCode := gjson.Get(json, "0")

}
func atom(tag string, kind string, functionCode uint16, address uint16, count uint16, period uint16, handler *modbus.TCPClientHandler) {
	client := modbus.NewClient(handler)
	switch functionCode {
	case variable.ReadCoil: //read coil
		for {
			results, _ := client.ReadCoils(uint16(address), count)
			log.Info("ReadCoil of ", tag, " is ", results)
			time.Sleep(5 * time.Second)
		}
	case variable.ReadDiscreteInput:
		for {
			results, _ := client.ReadDiscreteInputs(uint16(address), uint16(count))
			log.Info("ReadDiscreteInput of ", tag, " is ", results)
			time.Sleep(5 * time.Second)

		}
	case variable.ReadMultipleHoldingRegister:
		for {
			results, _ := client.ReadHoldingRegisters(uint16(address), uint16(count))
			log.Notice("ReadMultipleHoldingRegister of ", tag, " is ", results)
			time.Sleep(5 * time.Second)

		}
	case variable.ReadInputRegister:
		for {
			results, _ := client.ReadInputRegisters(uint16(address), uint16(count))
			log.Info("ReadInputRegister of ", tag, " is ", results)
			time.Sleep(5 * time.Second)

		}
	default:
		panic("fuuuuck")

	}

}
