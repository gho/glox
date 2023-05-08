package main

import "fmt"

type ValueType byte

const (
	ValueNil ValueType = iota
	ValueBool
	ValueNumber
)

type Value struct {
	typ  ValueType
	data interface{}
}

func nilValue() Value {
	return Value{typ: ValueNil}
}

func boolValue(b bool) Value {
	return Value{typ: ValueBool, data: b}
}

func numberValue(f float64) Value {
	return Value{typ: ValueNumber, data: f}
}

func (v Value) String() string {
	switch v.typ {
	case ValueNil:
		return "nil"
	case ValueBool:
		if v.data.(bool) {
			return "true"
		} else {
			return "false"
		}
	case ValueNumber:
		return fmt.Sprintf("%f", v.data)
	default:
		return "<unknown type>"
	}
}

func (v Value) asBool() bool {
	switch v.typ {
	case ValueBool:
		return v.data.(bool)
	case ValueNil:
		return v.data != nil
	}
	return true
}

func (v Value) asNumber() float64 {
	return v.data.(float64)
}

func negateValue(v Value) (Value, error) {
	return numberValue(-v.asNumber()), nil
}

func notValue(v Value) (Value, error) {
	return boolValue(!v.asBool()), nil
}

func addValues(v, w Value) (Value, error) {
	if v.typ == ValueNumber && w.typ == ValueNumber {
		return numberValue(v.asNumber() + w.asNumber()), nil
	}
	return Value{}, fmt.Errorf("type mismatch")
}

func subtractValues(v, w Value) (Value, error) {
	if v.typ == ValueNumber && w.typ == ValueNumber {
		return numberValue(v.asNumber() - w.asNumber()), nil
	}
	return Value{}, fmt.Errorf("type mismatch")
}

func multiplyValues(v, w Value) (Value, error) {
	if v.typ == ValueNumber && w.typ == ValueNumber {
		return numberValue(v.asNumber() * w.asNumber()), nil
	}
	return Value{}, fmt.Errorf("type mismatch")
}

func divideValues(v, w Value) (Value, error) {
	if v.typ == ValueNumber && w.typ == ValueNumber {
		return numberValue(v.asNumber() / w.asNumber()), nil
	}
	return Value{}, fmt.Errorf("type mismatch")
}

func valuesEqual(v, w Value) (Value, error) {
	res := false

	if v.typ == w.typ {
		switch v.typ {
		case ValueNil:
			res = true
		case ValueBool:
			res = v.asBool() == w.asBool()
		case ValueNumber:
			res = v.asNumber() == w.asNumber()
		}
	}

	return boolValue(res), nil
}

func valueGreater(v, w Value) (Value, error) {
	if v.typ == ValueNumber && w.typ == ValueNumber {
		return boolValue(v.asNumber() > w.asNumber()), nil
	}
	return Value{}, fmt.Errorf("type mismatch")
}

func valueLess(v, w Value) (Value, error) {
	if v.typ == ValueNumber && w.typ == ValueNumber {
		return boolValue(v.asNumber() < w.asNumber()), nil
	}
	return Value{}, fmt.Errorf("type mismatch")
}
