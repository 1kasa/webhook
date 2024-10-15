package utils

import jsoniter "github.com/json-iterator/go"

func MarshalAny2String(v interface{}) string {
	str,err := jsoniter.MarshalToString(v)
	if err != nil {
		return "{}"
	}
	return str
}