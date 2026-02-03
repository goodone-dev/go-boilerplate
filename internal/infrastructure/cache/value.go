package cache

import (
	"encoding/json"
	"strconv"
)

type CacheValue string

func (v *CacheValue) ToInt() int {
	if v == nil {
		return 0
	}

	val, _ := strconv.Atoi(string(*v))
	return val
}

func (v *CacheValue) ToInt32() int32 {
	if v == nil {
		return 0
	}

	val, _ := strconv.ParseInt(string(*v), 10, 32)
	return int32(val)
}

func (v *CacheValue) ToInt64() int64 {
	if v == nil {
		return 0
	}

	val, _ := strconv.ParseInt(string(*v), 10, 64)
	return val
}

func (v *CacheValue) ToFloat32() float32 {
	if v == nil {
		return 0
	}

	val, _ := strconv.ParseFloat(string(*v), 32)
	return float32(val)
}

func (v *CacheValue) ToFloat64() float64 {
	if v == nil {
		return 0
	}

	val, _ := strconv.ParseFloat(string(*v), 64)
	return val
}

func (v *CacheValue) ToBool() bool {
	if v == nil {
		return false
	}

	val, _ := strconv.ParseBool(string(*v))
	return val
}

func (v *CacheValue) ToString() string {
	if v == nil {
		return ""
	}

	return string(*v)
}

func (v *CacheValue) ToBytes() []byte {
	if v == nil {
		return nil
	}

	return []byte(string(*v))
}

func (v *CacheValue) IsEmpty() bool {
	return v == nil || string(*v) == ""
}

func (v *CacheValue) Decode(obj any) error {
	if v == nil {
		return nil
	}

	err := json.Unmarshal([]byte(string(*v)), &obj)
	return err
}
