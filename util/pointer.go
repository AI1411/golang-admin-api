package util

import "time"

func BoolToPtr(b bool) *bool {
	return &b
}

func StringToPtr(s string) *string {
	return &s
}

func Int8ToPtr(i int8) *int8 {
	return &i
}

func Int32ToPtr(i int32) *int32 {
	return &i
}

func Int64ToPtr(i int64) *int64 {
	return &i
}

func Uint8ToPtr(i uint8) *uint8 {
	return &i
}

func Uint32ToPtr(i uint32) *uint32 {
	return &i
}

func Uint64ToPtr(i uint64) *uint64 {
	return &i
}

func TimeToPtr(t time.Time) *time.Time {
	return &t
}
