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
	WriteMultipleCoils          uint16 = 6
	WriteSingleRegister         uint16 = 7
	WriteMultipleRegisters      uint16 = 8
)
