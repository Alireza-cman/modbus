package variable

const (

	//Configuration file
	ConfigFile = "config.json"

	//function codes
	ReadCoil                    uint16 = 1
	ReadDiscreteInput           uint16 = 2
	ReadMultipleHoldingRegister uint16 = 3
	ReadInputRegister           uint16 = 4
	WriteSingleCoil             uint16 = 5
	WriteSingleRegister         uint16 = 6
	WriteMultipleCoils          uint16 = 15
	WriteMultipleRegisters      uint16 = 16

	//
	//device to gateway Topic
	D2G_Sensor           = "v1/sensors"           // -t  "v1/sensors" -m {"serialNumber":"SN-001", "model":"T1000", "temperature":36.6}
	D2G_Sensor_embed     = "v1/sensors/+/+"       //sensor/SN-004/temperature
	D2G_Connect          = "v1/sensors/connect"   // -t "v1/sensors/connect" -m '{"serialNumber":"SN-001"}'
	D2G_Connect_embed    = "v1/sensors/+/connect" // -t "v1/sensors/SN-001/connect" -m ''
	D2G_Disconnect       = "v1/sensors/disconnect"
	D2G_Disconnect_embed = "v1/sensors/+/disconnect"
	D2G_RPC              = "v1/sensors/+/request/+/+" // v1/sensors/deviceName/request/method/requestID
	//
)
