// Code generated by "stringer -type WordNetLink wordnetlink.go"; DO NOT EDIT.

package acrostic

import "fmt"

const _WordNetLink_name = "WNSynonymWNHypeWNHypoWNMprtWNHprtWNHmemWNMmemWNMsubWNHsubWNDmncWNDmtcWNDmnuWNDmtuWNDmnrWNDmtrWNInstWNHasiWNEntaWNCausWNAlsoWNAttrWNSimWNEnd"

var _WordNetLink_index = [...]uint8{0, 9, 15, 21, 27, 33, 39, 45, 51, 57, 63, 69, 75, 81, 87, 93, 99, 105, 111, 117, 123, 129, 134, 139}

func (i WordNetLink) String() string {
	if i < 0 || i >= WordNetLink(len(_WordNetLink_index)-1) {
		return fmt.Sprintf("WordNetLink(%d)", i)
	}
	return _WordNetLink_name[_WordNetLink_index[i]:_WordNetLink_index[i+1]]
}
