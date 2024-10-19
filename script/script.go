package script

import _ "embed"

var (
	//go:embed redis/interaction_incr_cnt.lua
	LuaInteractionIncrCnt string
	//go:embed redis/send_code.lua
	LuaSendCode string
	//go:embed redis/verify_code.lua
	LuaVerifyCode string
)
