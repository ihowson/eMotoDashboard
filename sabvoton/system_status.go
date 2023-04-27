package sabvoton

type SystemStatus struct {
	// initial: [
	//     23 4445 0 65528 0 13 25 130 19933 65496 2 4503 426 62962 0 0
	//     0 1 0 65144 54452 50202 52916 50937 47101 0]
	Status uint16 // 23 'running with flux weakening'

	// these seem to align with the other registers

	// these all match the Status registers in registers.go
	Unknown1  float64 // 4445; float? sits around 0.7 (batt voltage?)
	Unknown2  float64 // ??? pos/neg weaken current command?
	Unknown3  float64 // 65528 float? around 10.0 weaken current feedback?
	Unknown4  float64 // ??? pos/neg torque current command
	Unknown5  float64 // ??? pos/neg torque current feedback
	Unknown6  uint16  // 25 unknown controller temperature
	Unknown7  uint16  // 130 unknown motor temperature
	Unknown8  float64 // 19933 float 0..5 but 2.5 center? motor angle
	Unknown9  float64 // 65496 float motor speed
	Unknown10 uint16  // values too small for double float, maybe uint16 or single float, only a few discrete values -- hall status?
	Unknown11 float64 // 4503 float maybe throttle voltage
	Unknown12 uint16  // 424 426 mosfet status?

	// these are unknown
	Unknown13 float64 // 65962 float
	Unknown14 uint16  // 0 unknown error code? +++
	Unknown15 uint16  // 0 unknown

	Unknown16 uint16  // 0 unknown
	Unknown17 uint16  // 1 unknown
	Unknown18 uint16  // 0 unknown
	Unknown19 float64 // 65144 float
	Unknown20 float64 // 54452 float
	Unknown21 uint16  // 50202 , 50234 in fault +++
	Unknown22 float64 // 52916 float
	Unknown23 float64 // 50937 float
	Unknown24 float64 // 47101 float
	Unknown25 uint16  // 0 unknown
}

func toFloat(in uint16, precision FloatPrecision) float64 {
	return float64(in) / float64(precision)
}

func parseSystemStatus(in []uint16) SystemStatus {
	return SystemStatus{
		Status:    in[0],
		Unknown1:  toFloat(in[1], Double),
		Unknown2:  toFloat(in[2], Double),
		Unknown3:  toFloat(in[3], Double),
		Unknown4:  toFloat(in[4], Double),
		Unknown5:  toFloat(in[5], Double),
		Unknown6:  in[6],
		Unknown7:  in[7],
		Unknown8:  toFloat(in[8], Double),
		Unknown9:  toFloat(in[9], Double),
		Unknown10: in[10],
		Unknown11: toFloat(in[11], Double),
		Unknown12: in[12],
		Unknown13: toFloat(in[13], Double),
		Unknown14: in[14],
		Unknown15: in[15],
		Unknown16: in[16],
		Unknown17: in[17],
		Unknown18: in[18],
		Unknown19: toFloat(in[19], Double),
		Unknown20: toFloat(in[20], Double),
		Unknown21: in[21],
		Unknown22: toFloat(in[22], Double),
		Unknown23: toFloat(in[23], Double),
		Unknown24: toFloat(in[24], Double),
		Unknown25: in[25],
	}
}
