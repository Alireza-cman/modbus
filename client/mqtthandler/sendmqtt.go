package mqtthandler

import (
	"encoding/hex"
	"strconv"

	"github.com/Atrovan/Modbus/client/variable"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	logging "github.com/op/go-logging"
	"github.com/tidwall/sjson"
)

var log = logging.MustGetLogger("example")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

var PublishClient MQTT.Client

func init() {
	opts := MQTT.NewClientOptions().AddBroker("tcp://localhost:1883")
	PublishClient = MQTT.NewClient(opts)
	if token := (PublishClient).Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Info("[DONE] The Modbus Application is connected to the Gateway")
}

func SendMQTTMessage(tag string, deviceName string, value []byte) {

	total, err := sjson.Set("", "serialNumber", deviceName)
	var onebit bool

	if err != nil {
		log.Error(err)
	}
	if len(value) == 0 {
		log.Error(" there is'nt value in there, are you sure your device is connected ? the value is empty")
		return
	} else if len(value) == 1 {
		onebit = true
	} else {
		onebit = false
	}
	if onebit == true {
		data := value[0]
		total, err = sjson.Set(total, tag, data)
		if token := PublishClient.Publish(variable.D2G_Sensor, 0, false, total); token.Wait() && token.Error() != nil {
			log.Error(token.Error())
		} else {
			log.Notice("coil has been read ")
		}
	} else {
		y := converter(value)
		total, err = sjson.Set(total, tag, y)
		if token := PublishClient.Publish(variable.D2G_Sensor, 0, false, total); token.Wait() && token.Error() != nil {
			log.Error(token.Error())
		} else {
			log.Notice("input register is working felan fine  ")
		}
	}
}
func bar(b []byte) string {
	return string(b)
}

func converter(input []byte) uint64 {
	//tmp := bar(input[:])
	y := hex.EncodeToString(input[:])
	n, err := strconv.ParseUint(y, 16, 32)
	if err != nil {
		log.Error(err)
		return 0
	} else {
		return n
	}
}
