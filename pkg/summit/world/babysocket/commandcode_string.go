// Code generated by "stringer -type=CommandCode"; DO NOT EDIT.

package babysocket

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[CommandPacket-0]
	_ = x[CommandInstruction-1]
	_ = x[CommandResponse-2]
}

const _CommandCode_name = "CommandPacketCommandInstructionCommandResponse"

var _CommandCode_index = [...]uint8{0, 13, 31, 46}

func (i CommandCode) String() string {
	if i >= CommandCode(len(_CommandCode_index)-1) {
		return "CommandCode(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _CommandCode_name[_CommandCode_index[i]:_CommandCode_index[i+1]]
}
