package helper

import (
	"strconv"
	"encoding/json"

	"xqdfs/errors"
)

func GetInt(param interface{}) (int,error){
	switch param.(type) {
	case string:
		v, err := strconv.Atoi(param.(string))
		if err == nil {
			return int(v), err
		} else {
			return 0, err
		}
	case json.Number:
		v, err := strconv.Atoi(param.(json.Number).String())
		if err == nil {
			return int(v), err
		} else {
			return 0, err
		}
	default:
		return 0,errors.ErrParameterError
	}
}

func GetInt32(param interface{}) (int32,error){
	switch param.(type) {
	case string:
		v, err := strconv.ParseInt(param.(string), 10, 32)
		if err == nil {
			return int32(v), err
		} else {
			return 0, err
		}
	case json.Number:
		v, err := strconv.ParseInt(param.(json.Number).String(), 10, 32)
		if err == nil {
			return int32(v), err
		} else {
			return 0, err
		}
	default:
		return 0,errors.ErrParameterError
	}
}

func GetInt64(param interface{}) (int64,error){
	switch param.(type) {
	case string:
		v, err := strconv.ParseInt(param.(string), 10, 64)
		if err == nil {
			return int64(v), err
		} else {
			return 0, err
		}
	case json.Number:
		v, err := strconv.ParseInt(param.(json.Number).String(), 10, 64)
		if err == nil {
			return int64(v), err
		} else {
			return 0, err
		}
	default:
		return 0,errors.ErrParameterError
	}
}
