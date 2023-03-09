package gui

import (
	"fmt"
	"log"
	"math"
	"time"

	imgui "github.com/inkyblackness/imgui-go/v4"
)

func (gui *MotoGUI) drawRunning() {
	DPMSForceOn()

	// FIXME: revamp
	m := gui.Model
	m.Lock.Lock()
	defer m.Lock.Unlock()

	/*
		top-right clock, date
		wifi symbol

		outside temperature

		big speedo in the middle
			small mph to its right
			- big circle guage around it from min to max speed like https://www.google.com/imgres?imgurl=https%3A%2F%2Fthumbor.forbes.com%2Fthumbor%2Ftrim%2F0x0%3A4000x2667%2Ffit-in%2F711x474%2Fsmart%2Fhttps%3A%2F%2Fspecials-images.forbesimg.com%2Fimageserve%2F5f488eecf326a401e79b743e%2FZero-SR-S-Electric-Motorcycle%2F0x0.jpg&imgrefurl=https%3A%2F%2Fwww.forbes.com%2Fsites%2Fbillroberson%2F2020%2F08%2F28%2Flong-term-ride-review-zeros-srs-electric-motorcycle-raises-the-bar-again%2F&tbnid=yBlBXHRIaz6XMM&vet=10CAMQxiAoAGoXChMIsNvgo9_k9wIVAAAAAB0AAAAAEBw..i&docid=plXQ59Hj1tpwiM&w=711&h=474&itg=1&q=electric%20motorcycle%20dashboard&ved=0CAMQxiAoAGoXChMIsNvgo9_k9wIVAAAAAB0AAAAAEBw
				- color can reflect something (% of max speed? temperature?)

		battery voltage graph on left with colors
		actual voltage
		% SOC
		est. range from last 10 miles
		est. range since start of charge

		lifetime odometer
		trip odometer

		motor temperature graph (scaled extenral temp to max temp)

		??
			current draw
			watts draw
				and a big pretty graph for that
			an rpm/speed graph might make more sense

		signal for left turn, right turn, highbeam, lights on
		any error indications

		flash battery when there's a range warning

		keep the background perfectly black, try to just have numbers floating on there

		when preset is 1, show N
		when preset is 2 or 3, show 'ECO' or 'SPORT'
		show the error flags reported by CA
	*/

	// if imgui.Button("Button") { // Buttons return true when clicked (most widgets return true when edited/activated)
	// 	counter++
	// }

	imgui.PushFont(gui.fontDINEng32)

	// STATUS ROW
	// speed graph & clock

	// Power graph.
	// On previous bikes this was kind of obvious; power was proportional to
	// throttle, and I know throttle from my hand position. It was also kind
	// of on/off, since there was little reason NOT to be at full throttle.
	// This bike is powerful enough that you will sit midrange for much of the
	// time.
	imgui.SetCursorPos(imgui.Vec2{
		X: 0.0,
		Y: 0.0,
	})
	imgui.ProgressBarV(
		float32(m.BatteryAmps/150.0),
		imgui.Vec2{
			X: 800.0,
			Y: 60.0,
		},
		fmt.Sprintf("%0.1fkW", m.BatteryVoltageCA*m.BatteryAmps/1000.0),
	)

	// TODO: status icons would go here

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

	// MIDDLE 3COL

	// use tables https://github.com/ocornut/imgui/blob/f58bd817e2998489bace1b2dff49884eac790efb/imgui_tables.cpp
	// then SetColumnWidth

	imgui.PushFont(gui.fontDINMittel32)

	col := imgui.StyleColorID(imgui.StyleColorText)

	color := imgui.Vec4{
		X: 0.5,
		Y: 1.0,
		Z: 0.5,
		W: 1.0,
	}

	// SPEEDO
	// TODO: push this as high as possible to help ergonomics
	imgui.SetCursorPos(imgui.Vec2{
		X: 240.0,
		Y: 120.0,
	})
	imgui.PushStyleColor(col, color)
	imgui.PushFont(gui.fontSpeed)
	// speed := m.LockNLoadFloat64(&m.Speed)
	intSpeed := int(math.Round(m.SpeedMph)) // TODO: use LockNLoad or an atomic read here

	// TODO: want this right-aligned. seems to 9-pad it with %2d
	strSpeed := fmt.Sprintf("%d", intSpeed)
	// %2d doesn't seem to pad spaces if the value is zero
	// also ' ' seems to be half-width? We should be locking the font width.
	// if intSpeed == 0 {
	// 	strSpeed = "  0"
	// }
	// strSpeed = ".00."
	imgui.Text(strSpeed)
	imgui.PopFont()
	imgui.PopStyleColor()

	imgui.SameLine()
	imgui.Text("mph")

	// battery bar
	imgui.PushFont(gui.fontDINEng32)
	// battWidth := float32(80.)
	battHeight := float32(gui.height)

	// might as well use all of the space available
	imgui.SetCursorPos(imgui.Vec2{
		X: 0.0,
		Y: 180.0,
	})

	// TODO: you should leave the X and Y width at 0.0f to have it automatically fit in the column
	fakeStateOfCharge := 1.0 - (m.BatteryAmpHoursConsumed / 16.0)
	VerticalProgressBar(
		float32(fakeStateOfCharge),
		imgui.Vec2{
			X: 80.0,
			Y: battHeight - 180.0,
		},
		fmt.Sprintf("bat\n%2d%%", int(math.Round(fakeStateOfCharge*100.0))),
	)

	// temperature gauges
	imgui.SetCursorPos(imgui.Vec2{
		X: 700.0,
		Y: 180.0,
	})
	VerticalProgressBar(
		float32(m.MotorTemperatureCelcius/100.0),
		imgui.Vec2{
			X: 80.0,
			Y: battHeight - 180.0,
		},
		fmt.Sprintf("mtr\n%d C", int(math.Round(m.MotorTemperatureCelcius))),
	)

	imgui.SetCursorPos(imgui.Vec2{
		X: 520.0,
		Y: 440.0,
	})
	imgui.Text(fmt.Sprintf("Trip: %0.1fmi", m.Distance))

	// imgui.SetCursorPos(imgui.Vec2{
	// 	X: 160.0,
	// 	Y: 280.0,
	// })
	// imgui.Text(fmt.Sprintf("Odometer: %0.0f", m.Odometer))

	imgui.SetCursorPos(imgui.Vec2{
		X: 120.0,
		Y: 440.0,
	})
	imgui.Text(fmt.Sprintf("BatteryVolts: %0.1f", m.BatteryVoltageCA))

	// TODO: report the minimum cell voltage here
	// - would like to see minimum battery voltage on the emoto on the main dashboard so you know if you're going to get a surprise low voltage; sparkline chart might be nice too

	imgui.SetCursorPos(imgui.Vec2{
		X: 120.0,
		Y: 400.0,
	})
	imgui.Text(fmt.Sprintf("Faults: %v", m.Faults))

	// Print Debugs
	// imgui.PushFont(imgui.DefaultFont)

	imgui.SetCursorPos(imgui.Vec2{
		X: 480.0,
		Y: 160.0,
	})
	imgui.Text("DEBUGS")
	imgui.SetCursorPos(imgui.Vec2{
		X: 480.0,
		Y: 180.0,
	})
	keys := SortedKeys(m.Debugs)
	y := float32(220.0)
	imgui.PushFont(imgui.DefaultFont)
	for _, key := range keys {
		log.Printf("try key=%s", key)
		value, ok := m.Debugs.Load(key)
		// log.Printf("Load key=%s, ok=%v value=%v", key, ok, value)
		if !ok {
			continue
		}
		// valueStr, ok := value.(string)
		// if !ok {
		// 	continue
		// }
		valueStr := fmt.Sprintf("%v", value)
		log.Printf("key=%s value=%s", key, valueStr)
		imgui.SetCursorPos(imgui.Vec2{X: 480.0, Y: y})
		imgui.Text(fmt.Sprintf("%s: %s", key, valueStr))
		y += 20.0
	}
	imgui.PopFont()

	// imgui.PushFont(gui.fontAwesome32)
	// imgui.Text("             ")

	// for i := 0xf200; i < 0xf300; i++ {
	// 	if i%32 == 0 {
	// 		imgui.Text("")
	// 	}
	// 	imgui.Text(fmt.Sprintf("%c", i))
	// 	imgui.SameLine()
	// }

	// imgui.Text(fmt.Sprintf("\uf175 %c %c %c", '', '', '\uf175')) // #           ")
	// speedo, wifi, temperature-half, bicycle, bolt, charging-station,
	// microchip (for controller temp?), plug, plug-circle-bolt,
	// plug-circle-exclamation, sliders, toggle-off, toggle-on, triangle-exclamation, motorcycle
	imgui.PopFont()
	imgui.PopFont()

	imgui.End()

	// TODO: insert skia/svg/whatever draw layer here

	// FIXME: maybe you can lock the size here?

	// FIXME: add a tap zone to go to Debugging or the other pages

}
