package gui

import imgui "github.com/inkyblackness/imgui-go/v4"

func (gui *MotoGUI) loadFonts(io imgui.IO) {
	// io.Fonts().AddFontDefault()
	gui.fontDINEng32 = io.Fonts().AddFontFromFileTTF("assets/DINEngschriftStd.otf", 32)
	gui.fontDINMittel32 = io.Fonts().AddFontFromFileTTF("assets/DINMittelschriftStd.otf", 32)

	// mph font
	config := imgui.NewFontConfig()
	config.SetGlyphMinAdvanceX(24.0)
	config.SetGlyphMaxAdvanceX(24.0)
	gui.fontSpeed = io.Fonts().AddFontFromFileTTF("assets/DINMittelschriftStd.otf", 240)

	// r := imgui.GlyphRangesBuilder{

	// ICON_MIN_FA := 0xf000
	// ICON_MAX_FA := 0xf2e0

	// imgui.New
	// range := imgui.NewGl

	// static const ImWchar icon_ranges[] = { ICON_MIN_FA, ICON_MAX_FA, 0 };
	// io.Fonts->AddFontFromFileTTF("fonts/fontawesome-webfont.ttf", 13.0f, &config, icon_ranges);

	// r := imgui.GlyphRanges(ICON_MIN_FA, ICON_MAX_FA)

	rb := imgui.GlyphRangesBuilder{}
	// rb.Add(ICON_MIN_FA'', '')
	rb.Add('\uf200', '\uf300')
	// rb.Add('\uf000', '\uf300')
	// rb.Add('\uf100', '\uf1ff')
	config.SetMergeMode(true)
	// rb.Add('\u2000', '\uf000')
	r := rb.Build()
	defer r.Free()
	gui.fontAwesome32 = io.Fonts().AddFontFromFileTTFV("assets/fa-regular-400.ttf", 32, config, r.GlyphRanges)

	_ = r
}
