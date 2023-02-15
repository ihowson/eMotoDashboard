package gpio

import "context"

type State struct {
	LeftBlinker  bool
	RightBlinker bool
	Headlights   bool
	HighBeam     bool
}

type GPIO interface {
	Run(ctx context.Context)
}
