package util

// derefString safely dereferences a pointer to a string.
func derefString(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}

// derefInt32 safely dereferences a pointer to an int32.
func derefInt32(ptr *int32) int32 {
	if ptr == nil {
		return 0
	}
	return *ptr
}

// derefBool safely dereferences a pointer to a bool.
func derefBool(ptr *bool) bool {
	if ptr == nil {
		return false
	}
	return *ptr
}
