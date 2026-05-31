class_name LoginPage
extends MarginContainer

signal connect_succeeded

@onready var login_panel: PanelContainer = $LoginCenter/LoginPanel
@onready var title_label: Label = $LoginCenter/LoginPanel/Margin/Body/LoginTitle
@onready var subtitle_label: Label = $LoginCenter/LoginPanel/Margin/Body/LoginSubtitle
@onready var endpoint_label: Label = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/EndpointField/EndpointLabel
@onready var protocol_label: Label = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/ProtocolField/ProtocolLabel
@onready var user_id_label: Label = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/UserIdField/UserIdLabel
@onready var token_label: Label = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/TokenField/TokenLabel
@onready var timeout_label: Label = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/TimeoutField/TimeoutLabel
@onready var reconnect_attempts_label: Label = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/ReconnectAttemptsField/ReconnectAttemptsLabel
@onready var reconnect_delay_label: Label = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/ReconnectDelayField/ReconnectDelayLabel
@onready var endpoint_input: LineEdit = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/EndpointField/EndpointInput
@onready var protocol_option: OptionButton = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/ProtocolField/ProtocolOption
@onready var user_id_input: LineEdit = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/UserIdField/UserIdInput
@onready var token_input: LineEdit = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/TokenField/TokenInput
@onready var timeout_input: LineEdit = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/TimeoutField/TimeoutInput
@onready var reconnect_attempts_input: LineEdit = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/ReconnectAttemptsField/ReconnectAttemptsInput
@onready var reconnect_delay_input: LineEdit = $LoginCenter/LoginPanel/Margin/Body/LoginGrid/ReconnectDelayField/ReconnectDelayInput
@onready var connect_button: Button = $LoginCenter/LoginPanel/Margin/Body/ConnectRow/ConnectButton
@onready var status_label: Label = $LoginCenter/LoginPanel/Margin/Body/LoginStatusLabel

var busy := false
var rpc_callback_target: Object = null


func _ready() -> void:
	_style_ui()
	_configure_protocol_options()
	connect_button.pressed.connect(_on_connect_pressed)
	_set_status("Ready.", UIStyle.TEXT_MUTED)
	_refresh_inputs()


func reset_after_disconnect(message: String = "Disconnected.", color: Color = UIStyle.TEXT_MUTED) -> void:
	busy = false
	_set_status(message, color)
	_refresh_inputs()


func _configure_protocol_options() -> void:
	protocol_option.clear()
	protocol_option.add_item("TCP", GolaxyClient.PROTOCOL_TCP)
	protocol_option.add_item("WebSocket", GolaxyClient.PROTOCOL_WEBSOCKET)
	protocol_option.select(1)


func _style_ui() -> void:
	UIStyle.style_panel(login_panel, UIStyle.CYAN)
	UIStyle.style_label(title_label, UIStyle.TEXT_PRIMARY, 34)
	UIStyle.style_label(subtitle_label, UIStyle.CYAN, 13)
	for label in [
		endpoint_label,
		protocol_label,
		user_id_label,
		token_label,
		timeout_label,
		reconnect_attempts_label,
		reconnect_delay_label,
	]:
		UIStyle.style_label(label, UIStyle.TEXT_MUTED, 12)

	for input in [
		endpoint_input,
		user_id_input,
		token_input,
		timeout_input,
		reconnect_attempts_input,
		reconnect_delay_input,
	]:
		UIStyle.style_line_edit(input)

	UIStyle.style_option_button(protocol_option)
	UIStyle.style_button(connect_button, UIStyle.CYAN)
	UIStyle.style_label(status_label, UIStyle.TEXT_MUTED)


func _on_connect_pressed() -> void:
	await _connect_to_server()


func _connect_to_server() -> void:
	var endpoint := endpoint_input.text.strip_edges()
	if endpoint.is_empty():
		_set_status("Endpoint is required.", UIStyle.AMBER)
		return

	busy = true
	_refresh_inputs()
	_set_status("Connecting to %s..." % endpoint, UIStyle.CYAN)
	if rpc_callback_target != null:
		RPCli.bind("", rpc_callback_target)

	var protocol_id := protocol_option.get_item_id(protocol_option.selected)
	var ok := await RPCli.connect_to_async(
		endpoint,
		protocol_id,
		user_id_input.text.strip_edges(),
		token_input.text.strip_edges(),
		_parse_int_or(timeout_input.text, GolaxyClient.DEFAULT_CONNECT_TIMEOUT_MS),
		_parse_int_or(reconnect_attempts_input.text, GolaxyClient.DEFAULT_RECONNECT_MAX_ATTEMPTS),
		_parse_int_or(reconnect_delay_input.text, GolaxyClient.DEFAULT_RECONNECT_DELAY_MS)
	)

	busy = false
	_refresh_inputs()
	if not ok:
		RPCli.unbind_script("")
		_set_status("Connect failed.", UIStyle.RED)
		return

	_set_status("Connected.", UIStyle.GREEN)
	connect_succeeded.emit()


func _refresh_inputs() -> void:
	connect_button.disabled = busy
	endpoint_input.editable = not busy
	protocol_option.disabled = busy
	user_id_input.editable = not busy
	token_input.editable = not busy
	timeout_input.editable = not busy
	reconnect_attempts_input.editable = not busy
	reconnect_delay_input.editable = not busy


func _set_status(message: String, color: Color) -> void:
	status_label.text = message
	status_label.add_theme_color_override("font_color", color)


func _parse_int_or(text: String, fallback: int) -> int:
	var value := text.strip_edges()
	if not value.is_valid_int():
		return fallback
	return int(value)
