package gui

import (
	"log"
	"os/exec"
)

var _timeout int

const (
	ForceOn  = -1
	ForceOff = -2
)

func DPMSForceOn() {
	if _timeout == ForceOn {
		return
	}
	_timeout = ForceOn
	cmd := exec.Command("xset", "s", "off", "-dpms")
	err := cmd.Run()
	if err != nil {
		log.Printf("Error setting DPMS timeout: %v", err)
	}
}

func DPMSForceOff() {
	if _timeout == ForceOff {
		return
	}
	_timeout = ForceOff
	cmd := exec.Command("xset", "dpms", "force", "off")
	err := cmd.Run()
	if err != nil {
		log.Printf("Error setting DPMS timeout: %v", err)
	}
}
