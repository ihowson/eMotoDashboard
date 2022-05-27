package gpio

type GPIOData struct {
	LeftBlinker  bool
	RightBlinker bool
	Headlights   bool
	HighBeam     bool
}

type GPIO interface {
	Run()
}
