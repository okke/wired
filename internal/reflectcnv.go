package internal

import (
	"reflect"
	"strconv"
)

type stringConvertor func(value string) reflect.Value

var stringConversions = make(map[reflect.Kind]stringConvertor, 0)

func init() {
	stringConversions[reflect.String] = string2string

	stringConversions[reflect.Int] = string2int
	stringConversions[reflect.Int8] = string2int8
	stringConversions[reflect.Int16] = string2int16
	stringConversions[reflect.Int32] = string2int32
	stringConversions[reflect.Int64] = string2int64

	stringConversions[reflect.Uint] = string2uint
	stringConversions[reflect.Uint8] = string2uint8
	stringConversions[reflect.Uint16] = string2uint16
	stringConversions[reflect.Uint32] = string2uint32
	stringConversions[reflect.Uint64] = string2uint64

	stringConversions[reflect.Float32] = string2float32
	stringConversions[reflect.Float64] = string2float64

	stringConversions[reflect.Bool] = string2bool
}

func string2string(value string) reflect.Value {
	return reflect.ValueOf(value)
}

func string2int(value string) reflect.Value {
	if i, err := strconv.Atoi(value); err == nil {
		return reflect.ValueOf(i)
	}
	return NilValue
}

func string2int8(value string) reflect.Value {
	if i, err := strconv.Atoi(value); err == nil {
		return reflect.ValueOf(int8(i))
	}
	return NilValue
}

func string2int16(value string) reflect.Value {
	if i, err := strconv.Atoi(value); err == nil {
		return reflect.ValueOf(int16(i))
	}
	return NilValue
}

func string2int32(value string) reflect.Value {
	if i, err := strconv.Atoi(value); err == nil {
		return reflect.ValueOf(int32(i))
	}
	return NilValue
}

func string2int64(value string) reflect.Value {
	if i, err := strconv.Atoi(value); err == nil {
		return reflect.ValueOf(int64(i))
	}
	return NilValue
}

func string2uint(value string) reflect.Value {
	if i, err := strconv.Atoi(value); err == nil {
		return reflect.ValueOf(uint(i))
	}
	return NilValue
}

func string2uint8(value string) reflect.Value {
	if i, err := strconv.Atoi(value); err == nil {
		return reflect.ValueOf(uint8(i))
	}
	return NilValue
}

func string2uint16(value string) reflect.Value {
	if i, err := strconv.Atoi(value); err == nil {
		return reflect.ValueOf(uint16(i))
	}
	return NilValue
}

func string2uint32(value string) reflect.Value {
	if i, err := strconv.Atoi(value); err == nil {
		return reflect.ValueOf(uint32(i))
	}
	return NilValue
}

func string2uint64(value string) reflect.Value {
	if i, err := strconv.Atoi(value); err == nil {
		return reflect.ValueOf(uint64(i))
	}
	return NilValue
}

func string2float32(value string) reflect.Value {
	if f, err := strconv.ParseFloat(value, 32); err == nil {
		return reflect.ValueOf(float32(f))
	}
	return NilValue
}

func string2float64(value string) reflect.Value {
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return reflect.ValueOf(f)
	}
	return NilValue
}

func string2bool(value string) reflect.Value {
	if f, err := strconv.ParseBool(value); err == nil {
		return reflect.ValueOf(f)
	}
	return NilValue
}

// ConvertString2Value takes a string and converts it to a value of given type
//
func ConvertString2Value(kind reflect.Kind, value string) reflect.Value {
	if converter, found := stringConversions[kind]; found {
		return converter(value)
	}
	return NilValue
}
