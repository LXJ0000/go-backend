package script

import _ "embed"

var (
	//go:embed redis/interaction_incr_cnt.lua
	LuaInteractionIncrCnt string
	//go:embed redis/slice_window.lua
	LuaSliceWindow string
)
