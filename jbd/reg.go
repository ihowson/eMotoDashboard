package jbd

const (
	RegisterEnterFactoryMode = 0x00
	RegisterExitFactoryMode  = 0x01
	// Register?? = 0x02
	RegisterBasicInfo    = 0x03
	RegisterCellVoltages = 0x04
	RegisterDeviceName   = 0x05
	RegisterUsePassword  = 0x06
	RegisterSetPassword  = 0x07
	// ?? 0x08?
	RegisterClearPassword = 0x09
)

type BasicInfoError uint16

const (
	BasicInfoErrorCellOverVolt         BasicInfoError = 0x0001
	BasicInfoErrorCellUnderVolt        BasicInfoError = 0x0002
	BasicInfoErrorPackOverVolt         BasicInfoError = 0x0004
	BasicInfoErrorPackUnderVolt        BasicInfoError = 0x0008
	BasicInfoErrorChargeOverTemp       BasicInfoError = 0x0010
	BasicInfoErrorChargeUnderTemp      BasicInfoError = 0x0020
	BasicInfoErrorDischargeOverTemp    BasicInfoError = 0x0040
	BasicInfoErrorDischargeUnderTemp   BasicInfoError = 0x0080
	BasicInfoErrorChargeOverCurrent    BasicInfoError = 0x0100
	BasicInfoErrorDischargeOverCurrent BasicInfoError = 0x0200
	BasicInfoErrorShortCircuit         BasicInfoError = 0x0400
	BasicInfoErrorFrontendICError      BasicInfoError = 0x0800
	BasicInfoErrorFETLocked            BasicInfoError = 0x1000
)
