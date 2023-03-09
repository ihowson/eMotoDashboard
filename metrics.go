package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/ihowson/eMotoDashboard/bike"
	"github.com/ihowson/eMotoDashboard/jbd"
)

// Runs a simple HTTP server to report metrics.

const tmpl = `
{{with .BasicInfo}}
PackVolts {{.PackVolts}}
PackAmps {{.PackAmps}}
PackCapacityAmpHours {{.PackCapacityAmpHours}} / {{.DesignCapacityAmpHours}}
StateOfChargePercent {{.StateOfChargePercent}}
InternalTemperature {{.InternalTemperature}}
PackTemperature1 {{.PackTemperature1}}
PackTemperature2 {{.PackTemperature2}}
{{end}}

{{with .CRate}}
C-Rate {{ printf "%0.2f" . }}C
{{end}}

{{range .CellVoltages.Volts}}
{{.}}V
{{end}}
`

type TemplateData struct {
	BasicInfo    *jbd.BasicInfo
	CellVoltages *jbd.CellVoltages
	CRate        float64
}

func MetricsServer(ctx context.Context, bike *bike.Bike) error {
	compiledTemplate := template.Must(template.New("metrics").Parse(tmpl))

	handler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, "metrics\n")

		basicInfo := bike.BMS.LatestBasicInfo()
		biPtr := &basicInfo
		var cRate float64
		if !(basicInfo.Time.IsZero() || time.Since(basicInfo.Time) < 5*time.Second) {
			biPtr = nil
		} else {
			cRate = basicInfo.PackAmps / 21.0 / 3.5 // 21 parallel cells, 3.5Ah per cell
		}

		cellVoltages := bike.BMS.LatestCellVoltages()
		cellVoltagesPtr := &cellVoltages
		if !(cellVoltages.Time.IsZero() || time.Since(cellVoltages.Time) < 5*time.Second) {
			cellVoltagesPtr = nil
		}

		log.Printf("cellVoltages %v", cellVoltages)

		data := TemplateData{
			BasicInfo:    biPtr,
			CellVoltages: cellVoltagesPtr,
			CRate:        cRate,
		}

		compiledTemplate.Execute(w, data)
	}

	// TODO: handle context

	http.HandleFunc("/", handler)
	return http.ListenAndServe(":8080", nil)
}
