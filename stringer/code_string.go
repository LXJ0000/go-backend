// Code generated by "stringer -type Code -linecomment"; DO NOT EDIT.

package stringer

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CODE_OK-0]
	_ = x[CODE_ERROR-1]
}

const _Code_name = "successfail"

var _Code_index = [...]uint8{0, 7, 11}

func (i Code) String() string {
	if i < 0 || i >= Code(len(_Code_index)-1) {
		return "Code(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Code_name[_Code_index[i]:_Code_index[i+1]]
}
