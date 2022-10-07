package lgr_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/nickbryan/collectable/libraries/lgr"
)

func TestField(t *testing.T) {
	t.Parallel()

	// Used to assert time for time fields.
	now := time.Now()

	testCases := map[string]struct {
		got, want lgr.Field
	}{
		"bool": {
			got: lgr.Bool("boolKey", true),
			want: lgr.Field{
				Type:  lgr.BoolType,
				Key:   "boolKey",
				Value: true,
			},
		},
		"byte string": {
			got: lgr.ByteStr("byteStrKey", []byte(`some byte string`)),
			want: lgr.Field{
				Type:  lgr.ByteStringType,
				Key:   "byteStrKey",
				Value: []byte(`some byte string`),
			},
		},
		"duration": {
			got: lgr.Duration("durationKey", time.Hour),
			want: lgr.Field{
				Type:  lgr.DurationType,
				Key:   "durationKey",
				Value: time.Hour,
			},
		},
		"error": {
			got: lgr.Err(errors.New("my error string")),
			want: lgr.Field{
				Type:  lgr.ErrorType,
				Key:   "error",
				Value: errors.New("my error string"),
			},
		},
		"named error": {
			got: lgr.NamedErr("namedErrKey", errors.New("my named error string")),
			want: lgr.Field{
				Type:  lgr.ErrorType,
				Key:   "namedErrKey",
				Value: errors.New("my named error string"),
			},
		},
		"float->float32": {
			got: lgr.Float("floatToFloat32Key", float32(4.20)),
			want: lgr.Field{
				Type:  lgr.Float32Type,
				Key:   "floatToFloat32Key",
				Value: float32(4.20),
			},
		},
		"float->float64": {
			got: lgr.Float("floatToFloat64Key", float64(0.24)),
			want: lgr.Field{
				Type:  lgr.Float64Type,
				Key:   "floatToFloat64Key",
				Value: float64(0.24),
			},
		},
		"integer->int": {
			got: lgr.Integer("integerToIntKey", int(42)),
			want: lgr.Field{
				Type:  lgr.IntType,
				Key:   "integerToIntKey",
				Value: int(42),
			},
		},
		"integer->int8": {
			got: lgr.Integer("integerToInt8Key", int8(42)),
			want: lgr.Field{
				Type:  lgr.Int8Type,
				Key:   "integerToInt8Key",
				Value: int8(42),
			},
		},
		"integer->int16": {
			got: lgr.Integer("integerToInt16Key", int16(42)),
			want: lgr.Field{
				Type:  lgr.Int16Type,
				Key:   "integerToInt16Key",
				Value: int16(42),
			},
		},
		"integer->int32": {
			got: lgr.Integer("integerToInt32Key", int32(42)),
			want: lgr.Field{
				Type:  lgr.Int32Type,
				Key:   "integerToInt32Key",
				Value: int32(42),
			},
		},
		"integer->int64": {
			got: lgr.Integer("integerToInt64Key", int64(42)),
			want: lgr.Field{
				Type:  lgr.Int64Type,
				Key:   "integerToInt64Key",
				Value: int64(42),
			},
		},
		"integer->uint": {
			got: lgr.Integer("integerToUintKey", uint(42)),
			want: lgr.Field{
				Type:  lgr.UintType,
				Key:   "integerToUintKey",
				Value: uint(42),
			},
		},
		"integer->uint8": {
			got: lgr.Integer("integerToUint8Key", uint8(42)),
			want: lgr.Field{
				Type:  lgr.Uint8Type,
				Key:   "integerToUint8Key",
				Value: uint8(42),
			},
		},
		"integer->uint16": {
			got: lgr.Integer("integerToUint16Key", uint16(42)),
			want: lgr.Field{
				Type:  lgr.Uint16Type,
				Key:   "integerToUint16Key",
				Value: uint16(42),
			},
		},
		"integer->uint32": {
			got: lgr.Integer("integerToUint32Key", uint32(42)),
			want: lgr.Field{
				Type:  lgr.Uint32Type,
				Key:   "integerToUint32Key",
				Value: uint32(42),
			},
		},
		"integer->uint64": {
			got: lgr.Integer("integerToUint64Key", uint64(42)),
			want: lgr.Field{
				Type:  lgr.Uint64Type,
				Key:   "integerToUint64Key",
				Value: uint64(42),
			},
		},
		"integer->uintptr": {
			got: lgr.Integer("integerToUintptrKey", uintptr(42)),
			want: lgr.Field{
				Type:  lgr.UintptrType,
				Key:   "integerToUintptrKey",
				Value: uintptr(42),
			},
		},
		"string": {
			got: lgr.Str("stringKey", "some string value"),
			want: lgr.Field{
				Type:  lgr.StringType,
				Key:   "stringKey",
				Value: "some string value",
			},
		},
		"time": {
			got: lgr.Time("timeKey", now),
			want: lgr.Field{
				Type:  lgr.TimeType,
				Key:   "timeKey",
				Value: now,
			},
		},
	}

	for testName, testCase := range testCases {
		tn, tc := testName, testCase

		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, tc.got)
		})
	}
}
