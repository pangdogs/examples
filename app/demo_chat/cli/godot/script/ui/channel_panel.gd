class_name ChannelPanel
extends PanelContainer

const Fnv32a = preload("res://addons/rpcli/fnv32a.gd")

signal create_channel_requested(channel: String)
signal join_channel_requested(channel: String)
signal switch_channel_requested(channel: String)
signal leave_channel_requested(channel: String)
signal remove_channel_requested(channel: String)
signal rtt_requested

@onready var title_label: Label = $Margin/Body/ChannelsTitle
@onready var active_channel_label: Label = $Margin/Body/ActiveChannelField/ActiveChannelFieldLabel
@onready var channel_input: LineEdit = $Margin/Body/ActiveChannelField/ChannelInput
@onready var create_button: Button = $Margin/Body/ChannelButtonGrid/CreateChannelButton
@onready var join_button: Button = $Margin/Body/ChannelButtonGrid/JoinChannelButton
@onready var switch_button: Button = $Margin/Body/ChannelButtonGrid/SwitchChannelButton
@onready var leave_button: Button = $Margin/Body/ChannelButtonGrid/LeaveChannelButton
@onready var remove_button: Button = $Margin/Body/ChannelButtonGrid/RemoveChannelButton
@onready var rtt_button: Button = $Margin/Body/ChannelButtonGrid/RttButton
@onready var channels_list: ItemList = $Margin/Body/ChannelsList
@onready var status_label: Label = $Margin/Body/ChatStatusLabel

var _buttons: Array[Button] = []


func _ready() -> void:
	_buttons.clear()
	_buttons.append(create_button)
	_buttons.append(join_button)
	_buttons.append(switch_button)
	_buttons.append(leave_button)
	_buttons.append(remove_button)
	_buttons.append(rtt_button)
	_style_ui()
	_connect_signals()
	set_status("", UIStyle.TEXT_MUTED)


func channel_text() -> String:
	return channel_input.text.strip_edges()


func set_channel(channel: String) -> void:
	channel_input.text = channel


func set_channels(channels: Array[String], current_channel: String) -> void:
	channels_list.clear()
	for channel in channels:
		var idx := channels_list.add_item(channel)
		channels_list.set_item_custom_fg_color(idx, _color_for(channel))
		if channel == current_channel:
			channels_list.select(idx)


func set_status(message: String, color: Color) -> void:
	status_label.text = message
	status_label.add_theme_color_override("font_color", color)


func set_enabled_state(connected: bool, busy: bool) -> void:
	channel_input.editable = connected and not busy
	channels_list.mouse_filter = Control.MOUSE_FILTER_STOP if connected and not busy else Control.MOUSE_FILTER_IGNORE
	for button in _buttons:
		button.disabled = busy or not connected


func _style_ui() -> void:
	UIStyle.style_panel(self, UIStyle.CYAN)
	UIStyle.style_label(title_label, UIStyle.TEXT_PRIMARY, 20)
	UIStyle.style_label(active_channel_label, UIStyle.TEXT_MUTED, 12)
	UIStyle.style_label(status_label, UIStyle.TEXT_MUTED)
	UIStyle.style_line_edit(channel_input)
	UIStyle.style_button(create_button, UIStyle.CYAN)
	UIStyle.style_button(join_button, UIStyle.GREEN)
	UIStyle.style_button(switch_button, UIStyle.MAGENTA)
	UIStyle.style_button(leave_button, UIStyle.AMBER)
	UIStyle.style_button(remove_button, UIStyle.RED)
	UIStyle.style_button(rtt_button, UIStyle.CYAN)


func _connect_signals() -> void:
	channel_input.text_submitted.connect(_on_channel_submitted)
	create_button.pressed.connect(func(): create_channel_requested.emit(channel_text()))
	join_button.pressed.connect(func(): join_channel_requested.emit(channel_text()))
	switch_button.pressed.connect(func(): switch_channel_requested.emit(channel_text()))
	leave_button.pressed.connect(func(): leave_channel_requested.emit(channel_text()))
	remove_button.pressed.connect(func(): remove_channel_requested.emit(channel_text()))
	rtt_button.pressed.connect(func(): rtt_requested.emit())
	channels_list.item_selected.connect(_on_channel_item_selected)
	channels_list.item_activated.connect(_on_channel_item_activated)


func _on_channel_submitted(_text: String) -> void:
	switch_channel_requested.emit(channel_text())


func _on_channel_item_selected(index: int) -> void:
	channel_input.text = channels_list.get_item_text(index)


func _on_channel_item_activated(index: int) -> void:
	channel_input.text = channels_list.get_item_text(index)
	switch_channel_requested.emit(channel_text())


func _color_for(value: String) -> Color:
	if value.is_empty():
		return UIStyle.TEXT_MUTED
	var hash := Fnv32a.new()
	hash.write_bytes(value.to_utf8_buffer())
	var h := hash.sum32()
	return Color.from_hsv(float(h % 360) / 360.0, 0.72, 0.94)
