package gui

import imgui "github.com/inkyblackness/imgui-go/v4"

func (gui *MotoGUI) loadFonts(io imgui.IO) {
	io.Fonts().AddFontDefault()
	gui.fontDINEng32 = io.Fonts().AddFontFromFileTTF("assets/DINEngschriftStd.otf", 32)
	gui.fontDINMittel32 = io.Fonts().AddFontFromFileTTF("assets/DINMittelschriftStd.otf", 32)

	// mph font
	// On the RPi3B+, I get blocked test if the font size is too big. It appears
	// that we're using too much GPU RAM or exceeding the maximum texture size.
	// To work around this, I only load the 0-9 glyphs and use minimal
	// oversampling.
	fontConfig := imgui.NewFontConfig()
	// Reduces texture size and quality. 2 is better.
	fontConfig.SetOversampleH(1)
	fontConfig.SetOversampleV(1)
	speedGlyphRangeBuilder := imgui.GlyphRangesBuilder{}
	speedGlyphRangeBuilder.Add('0', '9')
	speedGlyphRangeBuilder.Add('%', '%')
	speedGlyphRange := speedGlyphRangeBuilder.Build()
	// defer speedGlyphRange.Free()
	gui.fontSpeed = io.Fonts().AddFontFromFileTTFV("assets/DINMittelschriftStd.otf", 360, fontConfig, speedGlyphRange.GlyphRanges)

	// static const ImWchar icon_ranges[] = { ICON_MIN_FA, ICON_MAX_FA, 0 };
	// io.Fonts->AddFontFromFileTTF("fonts/fontawesome-webfont.ttf", 13.0f, &config, icon_ranges);

	// r := imgui.GlyphRanges(ICON_MIN_FA, ICON_MAX_FA)

	// rb := imgui.GlyphRangesBuilder{}
	// rb.Add(ICON_MIN_FA'', '')
	// TODO: you can add individual glyphs
	// config.SetMergeMode(true)
	// rb.Add('\u2000', '\uf000')
	// r := rb.Build()
	// defer r.Free()
	// gui.fontAwesome32 = io.Fonts().AddFontFromFileTTFV("assets/fa-regular-400.ttf", 32, config, r.GlyphRanges)
}
