package gui

import (
	"log"
	"os/exec"
	"strconv"
)

var _timeout int

func DPMSForceOn() {
	if 0 == _timeout {
		return
	}
	_timeout = 0
	cmd := exec.Command("xset", "dpms", "force", "on")
	err := cmd.Run()
	if err != nil {
		log.Printf("Error setting DPMS timeout: %v", err)
	}
}

func DPMSSetTimeout(seconds int) {
	if seconds == _timeout {
		return
	}
	_timeout = seconds

	strTimeout := strconv.Itoa(seconds)
	cmd := exec.Command("xset", "s", strTimeout, strTimeout)
	err := cmd.Run()
	if err != nil {
		log.Printf("Error setting DPMS timeout: %v", err)
	}
}
