package main

import (
	"github.com/ihowson/eMotoDashboard/m/v2/gui"
)

func main() {
	// TODO: disable screen blanking
	//     xset s off -dpms
	// then to turn it off programmatically
	//     xset dpms force off
	//     xset dpms force on

	model := BuildSystem()

	gui := gui.MotoGUI{
		Model: model,
	}
	gui.Run()

	// TODO: send cancel context into system and each component
	// cancel := ca.Run()
	// defer cancel()
}
