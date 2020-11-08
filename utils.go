package main

// Returns a pointer to the passed int32
func Int32Ptr(num int32) *int32 {
	return &num
}

// Returns a int value for the passed int32 pointer, zeros in case null
func Int32(num *int32) int32 {
	if num != nil {
		return *num
	}
	return 0
}