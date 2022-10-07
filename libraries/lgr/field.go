package lgr

import (
	"time"

	"golang.org/x/exp/constraints"
)

// FieldType indicates which member of the Field union struct should be used
// and how it should be serialised.
type FieldType uint8

const (
	// UnknownType is the default field type.
	UnkownType FieldType = iota
	// BoolType indicates that the field carries a boolean.
	BoolType
	// ByteStringType indicates that the field carries a UTF-8 encoding slice of bytes.
	ByteStringType
	// Duration indicates that the field carries a time.Duration.
	DurationType
	// ErrorType indicates that the field carries an error.
	ErrorType
	// Float32Type indicates that the field carries a float32.
	Float32Type
	// Float64Type indicates that the field carries a float64.
	Float64Type
	// IntType indicates that the field carries an int.
	IntType
	// Int8Type indicates that the field carries an int8.
	Int8Type
	// Int16Type indicates that the field carries an int16.
	Int16Type
	// Int32Type indicates that the field carries an int32.
	Int32Type
	// Int64Type indicates that the field carries an int64.
	Int64Type
	// UintType indicates that the field carries a uint.
	UintType
	// Uint8Type indicates that the field carries a u8int.
	Uint8Type
	// Uint16Type indicates that the field carries a uint16.
	Uint16Type
	// Uint32Type indicates that the field carries a uint32.
	Uint32Type
	// Uint64Type indicates that the field carries a uint64.
	Uint64Type
	// UintptrType indicates that the field carries a uintptr.
	UintptrType
	// StringType indicates that the field carries a string.
	StringType
	// TimeType indicates that the field carries a tiem.Time object.
	TimeType
)

// Field represents a key value pair that should be added to the log context. The
// Tyoe is used by the Adapter to write logs in a type safe way.
type Field struct {
	Type  FieldType
	Key   string
	Value any
}

// Bool constructs a Field that carries a boolean value with the given key.
func Bool(key string, value bool) Field {
	return Field{
		Type:  BoolType,
		Key:   key,
		Value: value,
	}
}

// ByteStr constructs a Field that carries a slice of UTF-8 encoded bytes with the given key.
func ByteStr(key string, value []byte) Field {
	return Field{
		Type:  ByteStringType,
		Key:   key,
		Value: value,
	}
}

// Duration constructs a Field that carries a time.Duration value with the given key.
func Duration(key string, value time.Duration) Field {
	return Field{
		Type:  DurationType,
		Key:   key,
		Value: value,
	}
}

// Err constructs a Field that carries an error value with the key "error".
func Err(err error) Field {
	return Field{
		Type:  ErrorType,
		Key:   "error",
		Value: err,
	}
}

// NamedErr constructs a Field that carries an error with the given key.
func NamedErr(key string, err error) Field {
	return Field{
		Type:  ErrorType,
		Key:   key,
		Value: err,
	}
}

// Float provides generic construction of a Field for float32 and float64 values.
func Float[T constraints.Float](key string, value T) Field {
	var typ FieldType

	switch any(value).(type) {
	case float32:
		typ = Float32Type
	case float64:
		typ = Float64Type
	}

	return Field{
		Type:  typ,
		Key:   key,
		Value: value,
	}
}

// Integer provides generic construction of a Field for all int and uint values.
func Integer[T constraints.Integer](key string, value T) Field {
	var typ FieldType

	switch any(value).(type) {
	case int:
		typ = IntType
	case int8:
		typ = Int8Type
	case int16:
		typ = Int16Type
	case int32:
		typ = Int32Type
	case int64:
		typ = Int64Type
	case uint:
		typ = UintType
	case uint8:
		typ = Uint8Type
	case uint16:
		typ = Uint16Type
	case uint32:
		typ = Uint32Type
	case uint64:
		typ = Uint64Type
	case uintptr:
		typ = UintptrType
	}

	return Field{
		Type:  typ,
		Key:   key,
		Value: value,
	}
}

// Str constructs a Field that carries a string value with the given key.
func Str(key, value string) Field {
	return Field{
		Type:  StringType,
		Key:   key,
		Value: value,
	}
}

// Time constructs a Field that carries a time.Time object with the given key.
func Time(key string, value time.Time) Field {
	return Field{
		Type:  TimeType,
		Key:   key,
		Value: value,
	}
}
