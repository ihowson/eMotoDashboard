package main

import (
	"github.com/ihowson/eMotoDashboard/bike"
	"github.com/ihowson/eMotoDashboard/gui"
)

func main() {
	// then to turn it off programmatically
	//     xset dpms force off
	//     xset dpms force on

	model, bike := bike.Build()

	gui := gui.MotoGUI{
		Model: model,
		Bike:  bike,
	}
	gui.Run()

	// TODO: send cancel context into system and each component
	// cancel := ca.Run()
	// defer cancel()
}
