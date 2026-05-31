class_name ChatPage
extends MarginContainer

signal signed_out(message: String, color: Color)

const GATE_SERVICE := "gate"
const CHAT_SERVICE := "chat"
const CHANNEL_COMP := "GateChatChannelComp"
const CHAT_USER_COMP := "ChatUserComp"
const GLOBAL_CHANNEL := "global"

@onready var header_panel: PanelContainer = $Page/HeaderPanel
@onready var title_label: Label = $Page/HeaderPanel/Margin/Body/HeaderRow/HeaderLeft/TitleRow/ChatTitle
@onready var online_badge: Label = $Page/HeaderPanel/Margin/Body/HeaderRow/HeaderLeft/TitleRow/OnlineBadge
@onready var active_channel_label: Label = $Page/HeaderPanel/Margin/Body/HeaderRow/HeaderLeft/ActiveChannelLabel
@onready var session_label: Label = $Page/HeaderPanel/Margin/Body/HeaderRow/HeaderLeft/SessionLabel
@onready var clock_label: Label = $Page/HeaderPanel/Margin/Body/HeaderRow/HeaderLeft/ClockLabel
@onready var disconnect_button: Button = $Page/HeaderPanel/Margin/Body/HeaderRow/DisconnectButton
@onready var channel_panel: ChannelPanel = $Page/Workspace/ChannelPanel
@onready var conversation_panel: ConversationPanel = $Page/Workspace/ConversationPanel

var current_channel := GLOBAL_CHANNEL
var connected := false
var busy := false
var _disconnecting_by_user := false
var _sending_message := false
var _probing_rtt := false
var known_channels: Array[String] = []


func _ready() -> void:
	_style_ui()
	_connect_ui_signals()
	_reset_channels()
	_set_connected(false)
	_set_channel(GLOBAL_CHANNEL)


func activate() -> void:
	_disconnecting_by_user = false
	busy = false
	_reset_channels()
	current_channel = GLOBAL_CHANNEL
	conversation_panel.clear_messages()

	var client := RPCli.client()
	if client != null and not client.disconnected.is_connected(_on_client_disconnected):
		client.disconnected.connect(_on_client_disconnected)

	_set_connected(true)
	_set_channel(GLOBAL_CHANNEL)
	conversation_panel.focus_input()
	_append_system_message("Connected. Global channel is active.")
	_update_clock_label(RPCli.remote_clock())
	_set_status("Connected.", UIStyle.GREEN)


func deactivate() -> void:
	RPCli.unbind("", get_instance_id())
	busy = false
	_set_connected(false)


func OutputText(timestamp: int, channel_name: String, id: String, text: String) -> void:
	var client := RPCli.client()
	var you := client != null and id == client.session_id()
	conversation_panel.append_message(timestamp, channel_name, id, text, you)


func ChannelKickOut(channel: String) -> void:
	_append_system_message("Kicked out from %s." % channel)
	_remove_known_channel(channel)
	if current_channel == channel:
		_set_channel(GLOBAL_CHANNEL)
		_set_status("Switched to global because %s is unavailable." % channel, UIStyle.AMBER)


func _style_ui() -> void:
	UIStyle.style_panel(header_panel, UIStyle.GREEN)
	UIStyle.style_label(title_label, UIStyle.TEXT_PRIMARY, 28)
	UIStyle.style_label(online_badge, UIStyle.GREEN, 13)
	UIStyle.style_label(active_channel_label, UIStyle.CYAN, 16)
	UIStyle.style_label(session_label, UIStyle.TEXT_MUTED)
	UIStyle.style_label(clock_label, UIStyle.MAGENTA)
	UIStyle.style_button(disconnect_button, UIStyle.RED)


func _connect_ui_signals() -> void:
	disconnect_button.pressed.connect(_on_disconnect_pressed)
	channel_panel.create_channel_requested.connect(_on_create_channel_requested)
	channel_panel.join_channel_requested.connect(_on_join_channel_requested)
	channel_panel.switch_channel_requested.connect(_on_switch_channel_requested)
	channel_panel.leave_channel_requested.connect(_on_leave_channel_requested)
	channel_panel.remove_channel_requested.connect(_on_remove_channel_requested)
	channel_panel.rtt_requested.connect(_on_rtt_requested)
	conversation_panel.send_requested.connect(_on_send_requested)


func _on_disconnect_pressed() -> void:
	if not connected:
		return
	_disconnecting_by_user = true
	RPCli.close(OK, "closed by user")
	_handle_signed_out("Disconnected.", UIStyle.TEXT_MUTED)


func _on_client_disconnected() -> void:
	if _disconnecting_by_user:
		return
	_handle_signed_out("Disconnected.", UIStyle.AMBER)


func _on_create_channel_requested(channel: String) -> void:
	var ok := await _channel_rpc("C_CreateChannel", channel)
	if ok:
		_add_known_channel(channel)
		_set_channel(channel)


func _on_remove_channel_requested(channel: String) -> void:
	var ok := await _channel_rpc("C_RemoveChannel", channel)
	if ok:
		_remove_known_channel(channel)
		if current_channel == channel:
			_set_channel(GLOBAL_CHANNEL)


func _on_join_channel_requested(channel: String) -> void:
	var ok := await _channel_rpc("C_JoinChannel", channel)
	if ok:
		_add_known_channel(channel)


func _on_leave_channel_requested(channel: String) -> void:
	if channel == GLOBAL_CHANNEL:
		_set_status("The global channel cannot be left.", UIStyle.AMBER)
		return
	var ok := await _channel_rpc("C_LeaveChannel", channel)
	if ok:
		_remove_known_channel(channel)
		if current_channel == channel:
			_set_channel(GLOBAL_CHANNEL)


func _on_switch_channel_requested(channel: String) -> void:
	if channel.is_empty():
		_set_status("Channel is required.", UIStyle.AMBER)
		return
	if channel == current_channel:
		return
	if channel == GLOBAL_CHANNEL:
		_set_channel(GLOBAL_CHANNEL)
		_set_status("Switched to global.", UIStyle.GREEN)
		return

	var result: Variant = await _rpc(GATE_SERVICE, CHANNEL_COMP, "C_InChannel", [channel])
	if result == null:
		return
	if bool(result.value):
		_add_known_channel(channel)
		_set_channel(channel)
		_set_status("Switched to %s." % channel, UIStyle.GREEN)
	else:
		_set_status("Join %s before switching to it." % channel, UIStyle.AMBER)


func _on_rtt_requested() -> void:
	if not _require_connected():
		return
	if _probing_rtt:
		return
	var client := RPCli.client()
	if client == null:
		return
	_probing_rtt = true
	var sample: GolaxyClient.TimeSample = await client.probe_time_async()
	_probing_rtt = false
	if sample == null:
		_set_status("RTT probe failed.", UIStyle.RED)
		return
	OutputText(int(Time.get_unix_time_from_system()), current_channel, client.session_id(), "RTT: %.3fs" % (float(sample.rtt()) / 1000.0))
	_update_clock_label(sample)


func _on_send_requested(text: String) -> void:
	if not _require_connected():
		return
	if _sending_message:
		return
	_sending_message = true
	var result: Variant = await _rpc(CHAT_SERVICE, CHAT_USER_COMP, "C_InputText", [current_channel, text], false)
	_sending_message = false
	if result == null:
		return
	conversation_panel.clear_input()


func _channel_rpc(method: String, channel: String) -> bool:
	if not _require_connected():
		return false
	if channel.is_empty():
		_set_status("Channel is required.", UIStyle.AMBER)
		return false
	if channel == GLOBAL_CHANNEL and method in ["C_CreateChannel", "C_RemoveChannel"]:
		_set_status("The global channel is managed by the server.", UIStyle.AMBER)
		return false

	var result: Variant = await _rpc(GATE_SERVICE, CHANNEL_COMP, method, [channel])
	if result == null:
		return false
	_set_status("%s %s succeeded." % [_display_rpc_method(method), channel], UIStyle.GREEN)
	return true


func _rpc(service: String, component: String, method: String, args: Array, update_busy: bool = true) -> Variant:
	if not _require_connected():
		return null
	if update_busy:
		_set_busy(true)
	var result: Variant = await RPCli.rpc_async(service, component, method, args)
	if update_busy:
		_set_busy(false)
	if result == null:
		_set_status("RPC %s.%s returned no result." % [component, method], UIStyle.RED)
		return null
	if not result.ok():
		_set_status("RPC %s.%s failed: %s" % [component, method, result.error], UIStyle.RED)
		return null
	return result


func _append_system_message(text: String) -> void:
	OutputText(int(Time.get_unix_time_from_system()), "system", "", text)


func _set_connected(value: bool) -> void:
	connected = value
	_refresh_enabled_state()
	if value:
		var client := RPCli.client()
		session_label.text = "Session: %s" % (client.session_id() if client != null else "-")
		online_badge.text = "ONLINE"
	else:
		session_label.text = "Session: none"
		online_badge.text = "OFFLINE"
		clock_label.text = "Clock: none"


func _set_busy(value: bool) -> void:
	busy = value
	_refresh_enabled_state()


func _refresh_enabled_state() -> void:
	disconnect_button.disabled = busy or not connected
	channel_panel.set_enabled_state(connected, busy)
	conversation_panel.set_enabled_state(connected, busy)


func _set_channel(channel: String) -> void:
	current_channel = channel
	_add_known_channel(channel)
	channel_panel.set_channel(channel)
	active_channel_label.text = "Channel: %s" % channel
	channel_panel.set_channels(known_channels, current_channel)


func _set_status(message: String, color: Color) -> void:
	channel_panel.set_status(message, color)


func _require_connected() -> bool:
	if connected:
		return true
	_set_status("Connect before sending RPC.", UIStyle.AMBER)
	return false


func _add_known_channel(channel: String) -> void:
	if channel.is_empty() or known_channels.has(channel):
		return
	known_channels.append(channel)
	channel_panel.set_channels(known_channels, current_channel)


func _reset_channels() -> void:
	known_channels.clear()
	known_channels.append(GLOBAL_CHANNEL)


func _remove_known_channel(channel: String) -> void:
	if channel == GLOBAL_CHANNEL:
		return
	known_channels.erase(channel)
	channel_panel.set_channels(known_channels, current_channel)


func _update_clock_label(sample: GolaxyClient.TimeSample) -> void:
	if sample == null:
		clock_label.text = "Clock: unknown"
		return
	var remote_time := Time.get_datetime_string_from_unix_time((sample.remote_now() / 1000) + sample.remote_zone_offset)
	clock_label.text = "Remote: %s | RTT: %dms" % [remote_time, sample.rtt()]


func _display_rpc_method(method: String) -> String:
	return method.substr(2) if method.begins_with("C_") else method


func _handle_signed_out(message: String, color: Color) -> void:
	if not connected:
		return
	_disconnecting_by_user = false
	busy = false
	RPCli.unbind("", get_instance_id())
	_set_connected(false)
	signed_out.emit(message, color)
