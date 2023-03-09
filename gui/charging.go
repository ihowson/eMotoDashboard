package gui

import (
	"fmt"
	"time"

	imgui "github.com/inkyblackness/imgui-go/v4"
)

func (gui *MotoGUI) drawCharging() {
	bike := gui.Bike
	bi := bike.BMS.LatestBasicInfo()
	biValid := bi.Time.IsZero() || time.Since(bi.Time) < 5*time.Second

	// gui.lock.Lock()
	// defer gui.lock.Unlock()

	/*
		battery voltage graph on left with colors
		actual voltage
	*/

	imgui.PushFont(gui.fontDINEng32)

	// STATUS ROW

	// clock
	now := time.Now()
	clockText := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())
	clockWidth := imgui.CalcTextSize(clockText, false, -1).X
	// imgui.SetCursorScreenPos(imgui.Vec2{
	imgui.SetCursorPos(imgui.Vec2{
		X: 800.0 - clockWidth - 5.0,
		Y: 65.0,
	})
	imgui.Text(clockText)
	imgui.PopFont()

	imgui.PushFont(gui.fontDINMittel32)

	// miles := range.EstimateMiles(bi)
	wattHoursPerMile := 60.0 // TOTAL GUESS
	miles := bi.PackCapacityAmpHours * 72.0 / wattHoursPerMile

	if biValid {
		col := imgui.StyleColorID(imgui.StyleColorText)
		green := imgui.Vec4{
			X: 0.0,
			Y: 1.0,
			Z: 0.0,
			W: 1.0,
		}
		red := imgui.Vec4{
			X: 1.0,
			Y: 0.0,
			Z: 0.0,
			W: 1.0,
		}

		imgui.Text(fmt.Sprintf("%d%% (%0.1fAh / %0.1fAh)", bi.StateOfChargePercent, bi.PackCapacityAmpHours, bi.DesignCapacityAmpHours))
		imgui.Text(fmt.Sprintf("%0.0f miles range", miles))

		// We would like to know if the charger is connected. This helps us
		// catch the case where the charger is connected but not charging.

		batteryFull := bi.PackVolts > (20 * 4.0)
		chargeRatekW := (bi.PackAmps * bi.PackVolts) / 1000.0
		charging := chargeRatekW > 0.0

		if batteryFull {
			imgui.PushStyleColor(col, green)
			imgui.Text("OK to disconnect charger.")
			imgui.PopStyleColor()
		}

		// TODO: PackAmps seems to be always positive -- or is it negative when discharging?
		if charging {
			DPMSForceOn()

			cRate := bi.PackAmps / 21.0 / 3.5 // 21 parallel cells, 3.5Ah per cell
			imgui.Text(fmt.Sprintf("%0.1fA / %0.1fkW / %0.2fC charge rate", bi.PackAmps, chargeRatekW, cRate))

			if !batteryFull {
				imgui.PushStyleColor(col, red)
				imgui.Text("Charging in progress. Do not disconnect charger.")
				imgui.PopStyleColor()
			}
		} else {
			imgui.Text("Not charging.")

			// you can go to sleep
			DPMSSetTimeout(60)
			time.Sleep(200 * time.Millisecond) // drag the chain, save some power
		}
	} else {
		imgui.Text("No BMS data")
	}

	imgui.Text(fmt.Sprintf("BMS: %0.1f°C", bi.InternalTemperature))
	imgui.Text(fmt.Sprintf("Battery: %0.1f°C %0.1f°C", bi.PackTemperature1, bi.PackTemperature2))

	cv := bike.BMS.LatestCellVoltages()
	cvValid := cv.Time.IsZero() || time.Since(cv.Time) < 5*time.Second
	// log.Printf("cv %v", cv.Volts)
	if cvValid && len(cv.Volts) > 0 {
		minV := cv.Volts[0]
		maxV := cv.Volts[0]
		for _, v := range cv.Volts {
			if v < minV {
				minV = v
			}
			if v > maxV {
				maxV = v
			}
		}
		imgui.Text(fmt.Sprintf("Cell min: %0.3fV max %0.3fV delta %0.1fmV", minV, maxV, (maxV-minV)*1000.0))
	}

	// 	// imgui.PlotHistogram("Cell Voltages", cv.Volts)
	// 	imgui.PlotHistogramV(
	// 		"Cell Voltages",
	// 		cv.Volts,
	// 		0,
	// 		"overlay",
	// 		0.0,
	// 		4.2, imgui.Vec2{X: 0, Y: 0})
	// }

	// imgui.Text("80% in 1h 20m") // time to 80% charge
	// imgui.Text("xx miles range")

	// col := imgui.StyleColorID(imgui.StyleColorText)

	// battery bar
	imgui.PushFont(gui.fontDINEng32)
	// battWidth := float32(80.)
	// battHeight := float32(gui.height)

	// might as well use all of the space available
	imgui.SetCursorPos(imgui.Vec2{
		X: 0.0,
		Y: 180.0,
	})

	// TODO: you should leave the X and Y width at 0.0f to have it automatically fit in the column
	// TODO: this is using the wrong battery capacity

	// temperature gauges
	/*
		imgui.SetCursorPos(imgui.Vec2{
			X: 700.0,
			Y: 180.0,
		})
		VerticalProgressBar(
			float32(m.BMSTemperatureCelcius/100.0),
			imgui.Vec2{
				X: 40.0,
				Y: battHeight - 180.0,
			},
			fmt.Sprintf("bms\n%d C", int(m.ControllerTemperatureCelcius)),
		)

		imgui.SetCursorPos(imgui.Vec2{
			X: 760.0,
			Y: 180.0,
		})
		VerticalProgressBar(
			float32(m.BatteryTemperature1Celcius/100.0),
			imgui.Vec2{
				X: 40.0,
				Y: battHeight - 180.0,
			},
			fmt.Sprintf("bat1\n%d C", int(math.Round(m.MotorTemperatureCelcius))),
		)

		imgui.SetCursorPos(imgui.Vec2{
			X: 760.0,
			Y: 180.0,
		})
		VerticalProgressBar(
			float32(m.BatteryTemperature1Celcius/100.0),
			imgui.Vec2{
				X: 40.0,
				Y: battHeight - 180.0,
			},
			fmt.Sprintf("bat2\n%d C", int(math.Round(m.MotorTemperatureCelcius))),
		)

		imgui.SetCursorPos(imgui.Vec2{
			X: 120.0,
			Y: 400.0,
		})
		imgui.Text(fmt.Sprintf("Faults: %v", m.Faults))

		// imgui.Text(fmt.Sprintf("\uf175 %c %c %c", '', '', '\uf175')) // #           ")
		// speedo, wifi, temperature-half, bicycle, bolt, charging-station,
		// microchip (for controller temp?), plug, plug-circle-bolt,
		// plug-circle-exclamation, sliders, toggle-off, toggle-on, triangle-exclamation, motorcycle
	*/
	imgui.PopFont()
	imgui.PopFont()

	imgui.End()
}
