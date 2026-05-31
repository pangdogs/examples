class_name UIStyle
extends RefCounted

const TEXT_PRIMARY := Color("f8fafc")
const TEXT_MUTED := Color("9ca3af")
const CYAN := Color("22d3ee")
const GREEN := Color("22c55e")
const AMBER := Color("f59e0b")
const RED := Color("ef4444")
const MAGENTA := Color("e879f9")
const BG_TOP := Color("070a12")
const BG_BOTTOM := Color("111827")
const PANEL_BG := Color(0.055, 0.074, 0.105, 0.88)
const FIELD_BG := Color(0.025, 0.035, 0.055, 0.92)
const BORDER := Color(0.24, 0.33, 0.45, 0.72)


static func create_theme(font_path: String) -> Theme:
	var font: Font = load(font_path) as Font
	if font == null:
		var system_font := SystemFont.new()
		system_font.font_names = PackedStringArray([
			"Microsoft YaHei",
			"SimHei",
			"Noto Sans CJK SC",
			"Noto Sans SC",
			"Source Han Sans SC",
			"PingFang SC",
			"WenQuanYi Micro Hei",
			"Arial Unicode MS",
			"sans-serif",
		])
		font = system_font

	var ui_theme := Theme.new()
	ui_theme.default_font = font
	return ui_theme


static func style_panel(panel: PanelContainer, accent: Color) -> void:
	var style := StyleBoxFlat.new()
	style.bg_color = PANEL_BG
	style.border_color = accent.darkened(0.15)
	style.border_width_left = 1
	style.border_width_top = 1
	style.border_width_right = 1
	style.border_width_bottom = 1
	style.corner_radius_top_left = 8
	style.corner_radius_top_right = 8
	style.corner_radius_bottom_left = 8
	style.corner_radius_bottom_right = 8
	style.shadow_color = Color(accent.r, accent.g, accent.b, 0.18)
	style.shadow_size = 18
	panel.add_theme_stylebox_override("panel", style)


static func style_label(label: Label, color: Color, font_size: int = 0) -> void:
	label.add_theme_color_override("font_color", color)
	if font_size > 0:
		label.add_theme_font_size_override("font_size", font_size)


static func style_line_edit(input: LineEdit) -> void:
	input.add_theme_color_override("font_color", TEXT_PRIMARY)
	input.add_theme_color_override("font_placeholder_color", TEXT_MUTED)
	input.add_theme_color_override("caret_color", CYAN)
	input.add_theme_stylebox_override("normal", field_style(BORDER))
	input.add_theme_stylebox_override("focus", field_style(CYAN))
	input.add_theme_stylebox_override("read_only", field_style(BORDER))


static func style_option_button(button: OptionButton) -> void:
	button.add_theme_color_override("font_color", TEXT_PRIMARY)
	button.add_theme_stylebox_override("normal", field_style(BORDER))
	button.add_theme_stylebox_override("focus", field_style(CYAN))
	button.add_theme_stylebox_override("hover", field_style(CYAN))


static func style_button(button: Button, accent: Color = CYAN) -> void:
	button.add_theme_color_override("font_color", TEXT_PRIMARY)
	button.add_theme_stylebox_override("normal", button_style(accent, false))
	button.add_theme_stylebox_override("hover", button_style(accent, true))
	button.add_theme_stylebox_override("pressed", button_style(accent.darkened(0.2), true))
	button.add_theme_stylebox_override("disabled", button_style(Color("374151"), false))


static func field_style(border_color: Color) -> StyleBoxFlat:
	var style := StyleBoxFlat.new()
	style.bg_color = FIELD_BG
	style.border_color = border_color
	style.border_width_left = 1
	style.border_width_top = 1
	style.border_width_right = 1
	style.border_width_bottom = 1
	style.corner_radius_top_left = 6
	style.corner_radius_top_right = 6
	style.corner_radius_bottom_left = 6
	style.corner_radius_bottom_right = 6
	style.content_margin_left = 10
	style.content_margin_right = 10
	return style


static func button_style(accent: Color, hover: bool) -> StyleBoxFlat:
	var style := StyleBoxFlat.new()
	style.bg_color = accent.darkened(0.45) if not hover else accent.darkened(0.28)
	style.border_color = accent
	style.border_width_left = 1
	style.border_width_top = 1
	style.border_width_right = 1
	style.border_width_bottom = 1
	style.corner_radius_top_left = 6
	style.corner_radius_top_right = 6
	style.corner_radius_bottom_left = 6
	style.corner_radius_bottom_right = 6
	style.shadow_color = Color(accent.r, accent.g, accent.b, 0.22 if hover else 0.10)
	style.shadow_size = 10 if hover else 4
	return style
