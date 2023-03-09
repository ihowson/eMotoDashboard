package gui

import (
	"time"
)

const timeout = 5 * time.Second

type State int

const (
	// Charging or idle.
	StateCharging = iota
	// Riding around.
	StateRunning
)

func (gui *MotoGUI) stateMachine() {
	// Every second, decide what state we should be in.
	bike := gui.Bike
	for {
		now := time.Now()

		// If the CA and Sabvoton are down, we might be idle or we might be charging.
		if now.After(bike.CycleAnalyst.LastUpdate().Add(timeout)) { //&&
			// now.After(bike.Sabvoton.LastUpdate().Add(timeout)) {
			gui.setState(StateCharging)
			// TODO: idle state
		} else {
			gui.setState(StateRunning)
		}

		time.Sleep(1 * time.Second)
	}
}

func (gui *MotoGUI) setState(state State) {
	gui.lock.Lock()
	defer gui.lock.Unlock()

	switch state {
	case StateCharging:
		gui.drawFunc = gui.drawCharging
	case StateRunning:
		gui.drawFunc = gui.drawRunning
	}
}
