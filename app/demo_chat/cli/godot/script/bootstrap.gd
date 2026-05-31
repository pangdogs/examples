extends Control

const CJK_FONT_PATH := "res://font/NotoSansCJKsc-Regular.otf"

@onready var login_page: LoginPage = $LoginPage
@onready var chat_page: ChatPage = $ChatPage


func _ready() -> void:
	set_anchors_and_offsets_preset(Control.PRESET_FULL_RECT)
	mouse_filter = Control.MOUSE_FILTER_PASS
	theme = UIStyle.create_theme(CJK_FONT_PATH)

	login_page.rpc_callback_target = chat_page
	login_page.connect_succeeded.connect(_on_login_connect_succeeded)
	chat_page.signed_out.connect(_on_chat_signed_out)
	_show_login()


func _exit_tree() -> void:
	RPCli.unbind_all()


func _on_login_connect_succeeded() -> void:
	login_page.visible = false
	chat_page.visible = true
	chat_page.activate()


func _on_chat_signed_out(message: String, color: Color) -> void:
	chat_page.deactivate()
	login_page.visible = true
	chat_page.visible = false
	login_page.reset_after_disconnect(message, color)


func _show_login() -> void:
	login_page.visible = true
	chat_page.visible = false
