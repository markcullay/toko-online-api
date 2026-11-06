package utils

import "strconv"

// ParseUintParam mengubah string ID dari URL menjadi tipe uint.
func ParseUintParam(s string) (uint, error) {
	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
