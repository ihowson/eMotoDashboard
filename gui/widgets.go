package gui

import (
	imgui "github.com/inkyblackness/imgui-go/v4"
)

func VerticalProgressBar(percent float32, size imgui.Vec2, label string) {
	style := imgui.CurrentStyle()
	bgColor := imgui.PackedColorFromVec4(style.Color(imgui.StyleColorTextSelectedBg)) // not quite the right color, but it'll do for now
	fgColor := imgui.PackedColorFromVec4(style.Color(imgui.StyleColorPlotHistogram))  // not quite the right color, but it'll do for now

	dl := imgui.WindowDrawList()
	bgMin := imgui.CursorPos()
	max := bgMin.Plus(size)

	fillHeight := size.Y * percent
	fgMin := imgui.Vec2{
		X: bgMin.X,
		Y: bgMin.Y + (size.Y - fillHeight),
	}

	// background fill
	dl.AddRectFilled(bgMin, max, bgColor)

	// foreground bar fill
	// log.Printf("fgMin=%v max=%v", fgMin, max)
	if max.Y > fgMin.Y {
		dl.AddRectFilled(fgMin, max, fgColor)
	} else {
		// log.Printf("out of range fgMin=%v max=%v", fgMin, max)
	}

	// center label in bar (ish) // FIXME: horizontal spacing is wrong
	textWidth := imgui.CalcTextSize(label, false, 0.0)
	midX := bgMin.X + (size.X / 2)
	startX := midX - textWidth.X/2
	imgui.SetCursorPos(imgui.Vec2{
		X: startX,
		Y: bgMin.Y + size.Y/2,
	})

	// TODO: use imgui.PushWrapPosV to center align
	imgui.Text(label)
}
