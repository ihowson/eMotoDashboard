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
	dl.AddRectFilled(fgMin, max, fgColor)

	// center label in bar (ish) // FIXME: horizontal spacing is wrong
	imgui.SetCursorPos(max.Minus(bgMin).Times(0.5).Plus(bgMin))

	// TODO: use imgui.PushWrapPosV to center align
	imgui.Text(label)
}
