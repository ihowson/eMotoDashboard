package gui

import (
	"fmt"
	"math"
	"time"

	imgui "github.com/inkyblackness/imgui-go/v4"

	"github.com/ihowson/eMotoDashboard/m/v2/dupe/platforms"
	"github.com/ihowson/eMotoDashboard/m/v2/dupe/renderers"
	"github.com/ihowson/eMotoDashboard/m/v2/model"
)

type MotoGUI struct {
	Model *model.Model

	fontDINEng32    imgui.Font
	fontDINMittel32 imgui.Font
	// fontAwesome32   imgui.Font
	fontSpeed imgui.Font

	clearColor [3]float32

	width  float32
	height float32
}

func (gui *MotoGUI) Run() error {
	// TODO: pull this from platform
	gui.width = 800.
	gui.height = 480.

	context := imgui.CreateContext(nil)
	defer context.Destroy()

	io := imgui.CurrentIO()
	gui.loadFonts(io)

	platform, err := platforms.NewGLFW(io, platforms.GLFWClientAPIOpenGL2)
	if err != nil {
		return fmt.Errorf("NewGLFW: %w", err)
	}
	defer platform.Dispose()

	renderer, err := renderers.NewOpenGL2(io)
	if err != nil {
		return fmt.Errorf("NewOpenGL2: %w", err)
	}
	defer renderer.Dispose()

	gui.run(platform, renderer)

	return nil
}

func (gui *MotoGUI) run(p Platform, r Renderer) {
	gui.clearColor = [3]float32{0.0, 0.0, 0.0}
	m := gui.Model

	for !p.ShouldStop() {
		p.ProcessEvents()

		// Signal start of a new frame
		p.NewFrame()
		imgui.NewFrame()

		gui.drawFrame(p, r, m)

		// ds := [2]float32{800.0, 480.0} // force screen size on Mac
		// r.Render(ds, p.FramebufferSize(), imgui.RenderedDrawData())
		r.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
		p.PostRender()
	}
}

func (gui *MotoGUI) drawFrame(p Platform, r Renderer, m *model.Model) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	// log.Printf(time.Now().String())

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

	imgui.SetNextWindowPos(imgui.Vec2{X: 0, Y: 0})
	imgui.SetNextWindowSize(imgui.Vec2{X: gui.width, Y: gui.height})
	imgui.BeginV("Dashboard", nil, imgui.WindowFlagsNoBackground|imgui.WindowFlagsNoDecoration)

	// if imgui.Button("Button") { // Buttons return true when clicked (most widgets return true when edited/activated)
	// 	counter++
	// }

	imgui.PushFont(gui.fontDINEng32)

	// STATUS ROW
	// speed graph & clock

	// RPM graph
	// actually speed, since we have no gears
	// TODO: want this to change color depending on speed/RPM
	imgui.SetCursorPos(imgui.Vec2{
		X: 0.0,
		Y: 0.0,
	})
	imgui.ProgressBarV(
		float32(m.SpeedMph/70.0),
		imgui.Vec2{
			X: 800.0,
			Y: 80.0,
		},
		"",
	)

	// TODO: status icons would go here

	// clock
	now := time.Now()
	clockText := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())
	clockWidth := imgui.CalcTextSize(clockText, false, -1).X
	// imgui.SetCursorScreenPos(imgui.Vec2{
	imgui.SetCursorPos(imgui.Vec2{
		X: 800.0 - clockWidth - 5.0,
		Y: 5.0,
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
	// speedoWidth := float32(400.)
	// speedoHeight := float32(300.)
	// imgui.SetCursorPos(imgui.Vec2{
	// 	X: (gui.width - speedoWidth) - (speedoWidth / 2),
	// 	Y: (gui.height - speedoHeight) - (speedoHeight / 2),
	// })
	imgui.SetCursorPos(imgui.Vec2{
		X: 240.0,
		Y: 120.0,
	})
	imgui.PushStyleColor(col, color)
	imgui.PushFont(gui.fontSpeed)

	// speed := m.LockNLoadFloat64(&m.Speed)
	intSpeed := int(math.Round(m.SpeedMph)) // TODO: use LockNLoad or an atomic read here
	// log.Printf("%d", intSpeed)
	// strSpeed := fmt.Sprintf("%02d", intSpeed)
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

	// imgui.WindowDrawList().AddCircleFilled(
	// 	imgui.Vec2{
	// 		X: 400,
	// 		Y: 240,
	// 	},
	// 	200.0,
	// 	imgui.PackedColor(0xaabbcc80),
	// )
	/*
		imgui.WindowDrawList().PathArcTo(
			imgui.Vec2{X: 400, Y: 240},
			100.0,
			0.25,
			0.75,
			12,
		)
	*/
	// func (list DrawList) PathArcTo(center Vec2, radius, a_min, a_max float32, num_segments int)

	// battery bar
	battWidth := float32(80.)
	battHeight := float32(gui.height)

	// might as well use all of the space available
	imgui.SetCursorPos(imgui.Vec2{
		X: 0.0,
		Y: 180.0,
	})

	// TODO: you should leave the X and Y width at 0.0f to have it automatically fit in the column
	fakeStateOfCharge := 1.0 - (m.BatteryAmpHoursConsumed / 16.0)
	imgui.ProgressBarV(
		// float32(m.BatteryStateOfCharge),
		float32(fakeStateOfCharge),
		imgui.Vec2{
			X: battWidth,
			Y: battHeight - 180.0,
		},
		// fmt.Sprintf("%2d%%", int(math.Round(m.BatteryStateOfCharge))),
		fmt.Sprintf("%2d%%", int(math.Round(fakeStateOfCharge*100.0))),
	)

	// temperature gaug
	// might as well use all of the space available
	imgui.SetCursorPos(imgui.Vec2{
		X: 720.0,
		Y: 180.0,
	})

	imgui.ProgressBarV(
		float32(m.MotorTemperatureCelcius/100.0),
		imgui.Vec2{
			X: battWidth,
			Y: battHeight - 180.0,
		},
		fmt.Sprintf("%d C", int(math.Round(m.MotorTemperatureCelcius))),
	)

	imgui.SetCursorPos(imgui.Vec2{
		X: 160.0,
		Y: 160.0,
	})
	imgui.Text(fmt.Sprintf("BatteryAmps: %0.1f", m.BatteryAmps))
	imgui.SetCursorPos(imgui.Vec2{
		X: 160.0,
		Y: 200.0,
	})
	imgui.Text(fmt.Sprintf("BatteryAmpHoursConsumed: %0.1f", m.BatteryAmpHoursConsumed))
	imgui.SetCursorPos(imgui.Vec2{
		X: 160.0,
		Y: 240.0,
	})
	imgui.Text(fmt.Sprintf("Distance: %0.1f", m.Distance))
	imgui.SetCursorPos(imgui.Vec2{
		X: 160.0,
		Y: 280.0,
	})
	imgui.Text(fmt.Sprintf("Odometer: %0.0f", m.Odometer))
	imgui.SetCursorPos(imgui.Vec2{
		X: 160.0,
		Y: 320.0,
	})
	imgui.Text(fmt.Sprintf("BatteryVolts: %0.1f", m.BatteryVoltageCA))
	imgui.SetCursorPos(imgui.Vec2{
		X: 160.0,
		Y: 360.0,
	})
	imgui.Text(fmt.Sprintf("Power: %0.2fkW", m.BatteryVoltageCA*m.BatteryAmps/1000.0))
	imgui.SetCursorPos(imgui.Vec2{
		X: 160.0,
		Y: 400.0,
	})
	imgui.Text(fmt.Sprintf("Gear: %s", m.Gear))
	imgui.SetCursorPos(imgui.Vec2{
		X: 160.0,
		Y: 440.0,
	})
	imgui.Text(fmt.Sprintf("Faults: %v", m.Faults))

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

	imgui.End()

	// Rendering
	imgui.Render() // This call only creates the draw data list. Actual rendering to framebuffer is done below.

	r.PreRender(gui.clearColor)
	// A this point, the application could perform its own rendering...
	// app.RenderScene()

	// TODO: insert skia/svg/whatever draw layer here

	// FIXME: maybe you can lock the size here?

	// FIXME: add a tap zone to go to Debugging or the other pages

}
