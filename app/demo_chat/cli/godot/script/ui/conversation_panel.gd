class_name ConversationPanel
extends VBoxContainer

const Fnv32a = preload("res://addons/rpcli/fnv32a.gd")

signal send_requested(text: String)

@onready var messages_panel: PanelContainer = $MessagesPanel
@onready var input_panel: PanelContainer = $InputPanel
@onready var messages_scroll: ScrollContainer = $MessagesPanel/Margin/Body/MessagesScroll
@onready var messages_box: VBoxContainer = $MessagesPanel/Margin/Body/MessagesScroll/MessagesBox
@onready var message_input: LineEdit = $InputPanel/Margin/Body/InputRow/MessageInput
@onready var send_button: Button = $InputPanel/Margin/Body/InputRow/SendButton

var max_messages := 256


func _ready() -> void:
	_style_ui()
	message_input.text_submitted.connect(_on_message_submitted)
	send_button.pressed.connect(_on_send_pressed)


func append_message(timestamp: int, channel_name: String, id: String, text: String, you: bool = false) -> void:
	var line := RichTextLabel.new()
	line.bbcode_enabled = true
	line.fit_content = true
	line.scroll_active = false
	line.selection_enabled = true
	line.size_flags_horizontal = Control.SIZE_EXPAND_FILL
	line.add_theme_color_override("default_color", UIStyle.TEXT_PRIMARY)
	line.add_theme_font_size_override("normal_font_size", 15)

	var user_text := "%s%s" % [id, " (YOU)" if you else ""]
	line.text = "%s  [color=%s]%s[/color]  [color=%s]%s[/color]: %s" % [
		_format_time(timestamp),
		_color_for(channel_name).to_html(false),
		_escape_bbcode(channel_name),
		_color_for(id).to_html(false),
		_escape_bbcode(user_text),
		_escape_bbcode(text),
	]
	messages_box.add_child(line)
	_trim_messages()
	_scroll_messages_to_bottom()


func clear_input() -> void:
	message_input.text = ""


func clear_messages() -> void:
	for child in messages_box.get_children():
		child.queue_free()


func focus_input() -> void:
	message_input.grab_focus()


func set_enabled_state(connected: bool, busy: bool) -> void:
	message_input.editable = connected and not busy
	send_button.disabled = busy or not connected


func _style_ui() -> void:
	UIStyle.style_panel(messages_panel, UIStyle.MAGENTA)
	UIStyle.style_panel(input_panel, UIStyle.GREEN)
	UIStyle.style_line_edit(message_input)
	UIStyle.style_button(send_button, UIStyle.GREEN)


func _on_send_pressed() -> void:
	_emit_send_requested()


func _on_message_submitted(_text: String) -> void:
	_emit_send_requested()


func _emit_send_requested() -> void:
	var text := message_input.text.strip_edges()
	if text.is_empty():
		return
	send_requested.emit(text)


func _trim_messages() -> void:
	while messages_box.get_child_count() > max_messages:
		messages_box.get_child(0).queue_free()


func _scroll_messages_to_bottom() -> void:
	await get_tree().process_frame
	if messages_scroll == null:
		return
	var bar := messages_scroll.get_v_scroll_bar()
	bar.value = bar.max_value


func _format_time(timestamp: int) -> String:
	var parts := Time.get_datetime_dict_from_unix_time(timestamp + _local_zone_offset_seconds())
	return "%02d:%02d:%02d" % [parts["hour"], parts["minute"], parts["second"]]


func _local_zone_offset_seconds() -> int:
	var zone := Time.get_time_zone_from_system()
	return int(zone.get("bias", 0)) * 60


func _color_for(value: String) -> Color:
	if value.is_empty():
		return UIStyle.TEXT_MUTED
	var hash := Fnv32a.new()
	hash.write_bytes(value.to_utf8_buffer())
	var h := hash.sum32()
	return Color.from_hsv(float(h % 360) / 360.0, 0.72, 0.94)


func _escape_bbcode(value: String) -> String:
	return value.replace("[", "<<LB>>").replace("]", "<<RB>>").replace("<<LB>>", "[lb]").replace("<<RB>>", "[rb]")
