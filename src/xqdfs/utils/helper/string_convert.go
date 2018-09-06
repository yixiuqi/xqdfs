package helper

import (
	"strconv"
	"errors"
)

func IntToString(value int) string {
	return strconv.Itoa(value)
}

func Int32ToString(value int32) string {
	var v int64 = int64(value)
	return strconv.FormatInt(v, 10)
}

func Int64ToString(value int64) string {
	return strconv.FormatInt(value, 10)
}

func StringToInt32(value string) (int32, error) {
	if value==""{
		return 0,errors.New("value is null")
	}

	v, err := strconv.ParseInt(value, 10, 32)
	if err == nil {
		return int32(v), err
	} else {
		return int32(0), err
	}
}

func StringToInt64(value string) (int64, error) {
	if value==""{
		return 0,errors.New("value is null")
	}

	v, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return int64(v), err
	} else {
		return int64(0), err
	}
}

func StringToUint64(value string) (uint64, error) {
	if value==""{
		return 0,errors.New("value is null")
	}

	return strconv.ParseUint(value, 10, 64)
}

func StringToInt(value string) (int, error) {
	if value==""{
		return 0,errors.New("value is null")
	}

	v, err := strconv.Atoi(value)
	if err == nil {
		return v, err
	} else {
		return 0, err
	}
}

func StringToFloat32(value string) (float32, error) {
	if value==""{
		return 0,errors.New("value is null")
	}

	v, err := strconv.ParseFloat(value, 32)
	if err == nil {
		return float32(v), err
	} else {
		return float32(0), err
	}
}

func StringToFloat64(value string) (float64, error) {
	if value==""{
		return 0,errors.New("value is null")
	}

	v, err := strconv.ParseFloat(value, 64)
	return v, err
}

func Float32ToString(value float32) string {
	return strconv.FormatFloat(float64(value), 'f', -1, 32)
}

func Float64ToString(value float64) string {
	return strconv.FormatFloat(value, 'f', -1, 64)
}
