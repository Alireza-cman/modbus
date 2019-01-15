package main

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/Atrovan/Modbus/client/rtuclient"
	"github.com/Atrovan/Modbus/client/tcpclient"
	"github.com/Atrovan/Modbus/client/variable"
	"github.com/op/go-logging"
	"github.com/tidwall/gjson"
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
		//fmt.Println(connection)
		if connection == "tcp" {
			go tcpclient.ModbusManipulationTCP(gjson.Get(configFile, "servers."+strconv.Itoa(i)).Raw)
		}
		if connection == "rtu" {
			go rtuclient.ModbusManipulationRTU(gjson.Get(configFile, "servers."+strconv.Itoa(i)).Raw)
		}

	}
	var wg sync.WaitGroup

	wg.Add(1)
	wg.Wait()

}
