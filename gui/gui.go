package gui

import (
	"fmt"
	"sync"

	"github.com/ihowson/eMotoDashboard/bike"
	"github.com/ihowson/eMotoDashboard/dupe/platforms"
	"github.com/ihowson/eMotoDashboard/dupe/renderers"
	"github.com/ihowson/eMotoDashboard/model"
	imgui "github.com/inkyblackness/imgui-go/v4"
)

type MotoGUI struct {
	Model *model.Model
	Bike  *bike.Bike

	lock sync.Mutex

	fontDINEng32    imgui.Font
	fontDINMittel32 imgui.Font
	// fontAwesome32   imgui.Font
	fontSpeed imgui.Font

	clearColor [3]float32

	width  float32
	height float32

	drawFunc func()

	platform Platform
	renderer Renderer
}

func (gui *MotoGUI) Run() error {
	go func() {
		gui.stateMachine()
	}()

	gui.drawFunc = gui.drawCharging // TODO: should be 'drawIdle'

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

	gui.platform = platform

	renderer, err := renderers.NewOpenGL2(io)
	if err != nil {
		return fmt.Errorf("NewOpenGL2: %w", err)
	}
	defer renderer.Dispose()

	gui.renderer = renderer

	gui.run()

	return nil
}

func (gui *MotoGUI) run() {
	gui.clearColor = [3]float32{0.0, 0.0, 0.0}

	p := gui.platform

	for !p.ShouldStop() {
		p.ProcessEvents()

		// Signal start of a new frame
		p.NewFrame()
		imgui.NewFrame()

		imgui.SetNextWindowPos(imgui.Vec2{X: 0, Y: 0})
		imgui.SetNextWindowSize(imgui.Vec2{X: gui.width, Y: gui.height})
		imgui.BeginV("Dashboard", nil, imgui.WindowFlagsNoBackground|imgui.WindowFlagsNoDecoration)

		gui.drawFunc()

		imgui.Render() // This call only creates the draw data list. Actual rendering to framebuffer is done below.
		gui.renderer.PreRender(gui.clearColor)

		// ds := [2]float32{800.0, 480.0} // force screen size on Mac
		gui.renderer.Render(p.DisplaySize(), p.FramebufferSize(), imgui.RenderedDrawData())
		p.PostRender()
	}
}
