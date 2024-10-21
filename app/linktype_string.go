// Code generated by "stringer -type=LinkType -linecomment"; DO NOT EDIT.

package app

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[LinkTypeUnspecified-0]
	_ = x[LinkTypeDocumentation-1]
	_ = x[LinkTypeAdministration-2]
	_ = x[LinkTypeSupport-3]
	_ = x[LinkTypeEndpoint-4]
	_ = x[LinkTypeExternal-5]
}

const _LinkType_name = "unspecifieddocumentationadministrationsupportendpointexternal"

var _LinkType_index = [...]uint8{0, 11, 24, 38, 45, 53, 61}

func (i LinkType) String() string {
	if i < 0 || i >= LinkType(len(_LinkType_index)-1) {
		return "LinkType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _LinkType_name[_LinkType_index[i]:_LinkType_index[i+1]]
}
