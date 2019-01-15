package rtuclient

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

//ModbusManipulationRTU is dealing with RTU connection
func ModbusManipulationRTU(configFile string) {

	//
	portName := gjson.Get(configFile, "transport.portName").Str
	encoding := gjson.Get(configFile, "transport.encoding").Str
	parity := gjson.Get(configFile, "transport.parity").Str
	timeout := gjson.Get(configFile, "transport.timeout").Num
	baudRate := gjson.Get(configFile, "transport.baudRate").Num
	dataBits := gjson.Get(configFile, "transport.dataBits").Num
	stopBits := gjson.Get(configFile, "transport.stopBits").Num

	if timeout == 0 {
		timeout = 0
	}
	if baudRate == 0 {
		baudRate = 19200
	}
	if encoding == "" {
		encoding = "rtu"
	}
	if parity == "" || parity == "None" || parity == "none" {
		parity = "N"
	}
	if portName == "" {
		log.Error("You didn't indicate the portName in the config file")
	}
	log.Critical("stopBits", stopBits)
	log.Critical("dataBits", dataBits)
	log.Critical("Baudrate", baudRate)
	log.Critical("parity", parity)
	log.Critical("encoding", encoding)
	log.Critical("portName", portName)

	handler := modbus.NewRTUClientHandler(portName)

	handler.BaudRate = int(baudRate)
	handler.Parity = parity
	handler.IdleTimeout = time.Duration(timeout) * time.Second
	handler.StopBits = int(stopBits)
	handler.DataBits = int(dataBits)
	handler.SlaveId = 1

	deviceNumber := gjson.Get(configFile, "devices.#").Num

	if encoding == "rtu" {
		log.Warning("trying to make a RTU modbus connection,", portName)
		for i := 0; i < int(deviceNumber); i++ {
			unitId := gjson.Get(configFile, "devices."+strconv.Itoa(i)+".unitId").Num
			deviceName := gjson.Get(configFile, "devices."+strconv.Itoa(i)+".deviceName").Str
			telemetry := gjson.Get(configFile, "devices."+strconv.Itoa(i)+".timeseries").Raw
			go TelemetryHandlerRTU(telemetry, unitId, deviceName, handler)

		}
	}

	//
}

func TelemetryHandlerRTU(json string, unitId float64, deviceName string, handler *modbus.RTUClientHandler) {
	telemetryNumber := gjson.Get(json, "#").Num
	log.Warning("telemetry number of device ", deviceName, "is:", telemetryNumber)
	for i := 0; i < int(telemetryNumber); i++ {
		tag := gjson.Get(json, strconv.Itoa(i)+".tag").Str
		kind := gjson.Get(json, strconv.Itoa(i)+".kind").Str
		functionCode := gjson.Get(json, strconv.Itoa(i)+".functionCode").Num
		address := gjson.Get(json, strconv.Itoa(i)+".address").Num
		count := gjson.Get(json, strconv.Itoa(i)+".registerCount").Num
		period := gjson.Get(json, strconv.Itoa(i)+".pollPeriod").Num
		if period == 0 {
			period = 10
		}
		if count == 0 {
			count = 2
		}
		go AtomRTU(tag, kind, uint16(functionCode), uint16(address), uint16(count), uint16(period), handler)

	}

	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()
	//functionCode := gjson.Get(json, "0")
}

func AtomRTU(tag string, kind string, functionCode uint16, address uint16, count uint16, period uint16, handler *modbus.RTUClientHandler) {
	err := handler.Connect()
	if err != nil {
		log.Error(err)
	}
	log.Critical(address)
	log.Critical(count)
	log.Critical(functionCode)
	client := modbus.NewClient(handler)
	//log.Critical(client)
	switch functionCode {
	case variable.ReadCoil: //read coil
		for {
			results, err := client.ReadCoils(uint16(address), count)
			if results[0] == byte(0) {
				_, err = client.WriteSingleCoil(address, uint16(65280))

			} else {
				_, err = client.WriteSingleCoil(address, uint16(0))
			}
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
		log.Error(functionCode)
		panic("fuuuuck")

	}

}
