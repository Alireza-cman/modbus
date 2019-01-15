package tcpclient

import (
	"strconv"
	"sync"
	"time"

	"github.com/Atrovan/Modbus/client/variable"
	"github.com/goburrow/modbus"
	logging "github.com/op/go-logging"
	"github.com/tidwall/gjson"
)

var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

// this function retrieve the host and port of the slave device, moreover, it retrieves the timeout which is necessary
// for the connection. at the end of the function, telemetry values which is needed for us will be fetched
func ModbusManipulationTCP(configFile string) {
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
		go TelemetryHandlerTCP(telemetry, unitId, deviceName, handler)
		//fmt.Println(telemetry)
		//fmt.Println(unitId, deviceName)

	}

	// var wg sync.WaitGroup

	// wg.Add(1)
	// wg.Wait()

}

func TelemetryHandlerTCP(json string, unitId float64, deviceName string, handler *modbus.TCPClientHandler) {
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
			period = 10
		}
		if count == 0 {
			count = 2
		}

		go AtomTCP(tag, kind, uint16(functionCode), uint16(address), uint16(count), uint16(period), handler)

	}

	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()
	//functionCode := gjson.Get(json, "0")
}

/*
Function Code	Register Type
1				Read Coil
2				Read Discrete Input
3				Read Holding Registers
4				Read Input Registers
5				Write Single Coil
6				Write Single Holding Register
15				Write Multiple Coils
16				Write Multiple Holding Registers
*/
func AtomTCP(tag string, kind string, functionCode uint16, address uint16, count uint16, period uint16, handler *modbus.TCPClientHandler) {
	err := handler.Connect()
	if err != nil {
		log.Error(err)
	}
	client := modbus.NewClient(handler)
	switch functionCode {
	case variable.ReadCoil: //read coil
		for {
			results, err := client.ReadCoils(uint16(address), count)
			//results, err := client.WriteSingleCoil(address, uint16(65280))
			if err != nil {
				log.Error(err)
			} else {
				log.Info("ReadCoil of ", tag, " is ", results)
			}

			time.Sleep(time.Duration(period) * time.Second)
		}
	case variable.ReadDiscreteInput:
		for {
			results, err := client.ReadDiscreteInputs(uint16(address), uint16(count))
			if err != nil {
				log.Error("ReadDiscreteInput of ", tag, " is ", results)
				log.Error(err)
			} else {
				log.Info("ReadDiscreteInput of ", tag, " is ", results)
			}
			time.Sleep(time.Duration(period) * time.Second)
		}
	case variable.ReadMultipleHoldingRegister:
		for {
			results, err := client.ReadHoldingRegisters(uint16(address), uint16(count))
			if err != nil {
				log.Error("ReadDiscreteInput of ", tag, " is ", results)
				log.Error(err)
			} else {
				log.Notice("ReadMultipleHoldingRegister of ", tag, " is ", results)

			}
			time.Sleep(time.Duration(period) * time.Second)

		}
	case variable.ReadInputRegister:
		for {
			results, err := client.ReadInputRegisters(uint16(address), uint16(count))
			if err != nil {
				log.Error("ReadDiscreteInput of ", tag, " is ", results)
				log.Error(err)
			} else {
				log.Info("ReadInputRegister of ", tag, " is ", results)
			}

			time.Sleep(time.Duration(period) * time.Second)
		}
	default:
		panic("fuuuuck")

	}

}
