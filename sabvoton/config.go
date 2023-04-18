package sabvoton

type Config struct {
	Register RegisterUInt16
	Name     string
	Value    uint16
}

var DesiredConfig = []Config{
	// Apparently you can change this and it's the DC limit cutout?
	{RegisterRatedDCCurrent, "Rated DC Current", 200},
	{RegisterMaxDCCurrent, "Max DC Current", 150},
	{RegisterMaxPhaseCurrent, "Max Phase Current", 280},
	// {RegisterProtectedPhaseCurrent, "Protected Phase Current", 450},
	// {RegisterRatedPhaseCurrent, "Rated Phase Current", 150},
	// {RegisterFluxWeakeningEnabled, "Flux Weakening Enabled", 0},
	// {RegisterFluxWeakenCurrent, "Flux Weaken Current", 0},
	// {RegisterThrottleMidVoltage, "Throttle Mid Voltage", 50},
	{RegisterThrottleMidPhaseCurrent, "Throttle Mid Phase Current", 100},
	// {RegisterThrottleMinVoltage, "Throttle Min Voltage", 0},
	// {RegisterThrottleMaxVoltage, "Throttle Max Voltage", 100},
}
