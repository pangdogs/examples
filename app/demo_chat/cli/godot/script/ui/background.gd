class_name Background
extends Control

var _fx_time := 0.0


func _ready() -> void:
	mouse_filter = Control.MOUSE_FILTER_IGNORE
	set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)


func _process(delta: float) -> void:
	_fx_time += delta
	queue_redraw()


func _draw() -> void:
	var rect := get_rect()
	draw_rect(rect, UIStyle.BG_TOP)
	var band_count := 10
	for i in range(band_count):
		var t := float(i) / float(maxi(1, band_count - 1))
		var y := rect.size.y * t
		var color := UIStyle.BG_TOP.lerp(UIStyle.BG_BOTTOM, t)
		draw_rect(Rect2(0, y, rect.size.x, rect.size.y / band_count + 1.0), color)

	var grid_color := Color(0.15, 0.23, 0.33, 0.20)
	var offset := fmod(_fx_time * 18.0, 48.0)
	for x in range(-48, int(rect.size.x) + 96, 48):
		draw_line(Vector2(x + offset, 0), Vector2(x - 180 + offset, rect.size.y), grid_color, 1.0)
	for y in range(-48, int(rect.size.y) + 96, 48):
		draw_line(Vector2(0, y + offset), Vector2(rect.size.x, y - 80 + offset), grid_color, 1.0)

	var p1 := Vector2(rect.size.x * 0.18, rect.size.y * 0.22)
	var p2 := Vector2(rect.size.x * 0.82, rect.size.y * 0.78)
	draw_circle(p1, 180.0 + sin(_fx_time) * 18.0, Color(0.03, 0.55, 0.68, 0.08))
	draw_circle(p2, 220.0 + cos(_fx_time * 0.8) * 22.0, Color(0.62, 0.22, 0.80, 0.07))
