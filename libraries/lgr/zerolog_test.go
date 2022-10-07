package lgr //nolint: testpackage // Integration test with zerolog.

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestZerologAdapter(t *testing.T) {
	t.Parallel()

	t.Run("panics when level is unexpected", func(t *testing.T) {
		t.Parallel()
		require.PanicsWithValue(t, "log level unexpected: 100", func() {
			adapter := zerologAdapter{logger: zerolog.New(&bytes.Buffer{})}
			adapter.Adapt(Level(100), "")
		})
	})

	// Used to assert time for time fields.
	now := time.Now()

	testCases := map[string]struct {
		level  Level
		msg    string
		fields []Field
		want   string
	}{
		"logs debug message": {
			level: DebugLevel,
			msg:   "my debug message",
			want:  `{"level":"debug","context":{},"message":"my debug message"}`,
		},
		"logs info message": {
			level: InfoLevel,
			msg:   "my info message",
			want:  `{"level":"info","context":{},"message":"my info message"}`,
		},
		"logs warn message": {
			level: WarnLevel,
			msg:   "my warn message",
			want:  `{"level":"warn","context":{},"message":"my warn message"}`,
		},
		"logs error message": {
			level: ErrorLevel,
			msg:   "my error message",
			want:  `{"level":"error","context":{},"message":"my error message"}`,
		},
		"sets bool field": {
			fields: []Field{{Type: BoolType, Key: "boolFieldKey", Value: true}},
			want:   `{"level":"info","context":{"boolFieldKey":true}}`,
		},
		"sets byte string field": {
			fields: []Field{{Type: ByteStringType, Key: "byteStringFieldKey", Value: []byte("some byte string")}},
			want:   `{"level":"info","context":{"byteStringFieldKey":"some byte string"}}`,
		},
		"sets duration field": {
			fields: []Field{{Type: DurationType, Key: "durationFieldKey", Value: time.Second}},
			want:   `{"level":"info","context":{"durationFieldKey":1000}}`,
		},
		"sets error field": {
			fields: []Field{{Type: ErrorType, Key: "errorFieldKey", Value: errors.New("some error string")}},
			want:   `{"level":"info","context":{"errorFieldKey":"some error string"}}`,
		},
		"sets float fields": {
			fields: []Field{
				{Type: Float32Type, Key: "float32FieldKey", Value: float32(123)},
				{Type: Float64Type, Key: "float64FieldKey", Value: float64(456)},
			},
			want: `{"level":"info","context":{"float32FieldKey":123,"float64FieldKey":456}}`,
		},
		"sets int fields": {
			fields: []Field{
				{Type: IntType, Key: "intFieldKey", Value: int(1)},
				{Type: Int8Type, Key: "int8FieldKey", Value: int8(2)},
				{Type: Int16Type, Key: "int16FieldKey", Value: int16(3)},
				{Type: Int32Type, Key: "int32FieldKey", Value: int32(4)},
				{Type: Int64Type, Key: "int64FieldKey", Value: int64(5)},
			},
			want: `{"level":"info","context":{"intFieldKey":1,"int8FieldKey":2,"int16FieldKey":3,"int32FieldKey":4,"int64FieldKey":5}}`,
		},
		"sets uint fields": {
			fields: []Field{
				{Type: UintType, Key: "uintFieldKey", Value: uint(1)},
				{Type: Uint8Type, Key: "uint8FieldKey", Value: uint8(2)},
				{Type: Uint16Type, Key: "uint16FieldKey", Value: uint16(3)},
				{Type: Uint32Type, Key: "uint32FieldKey", Value: uint32(4)},
				{Type: Uint64Type, Key: "uint64FieldKey", Value: uint64(5)},
				{Type: UintptrType, Key: "uintptrFieldKey", Value: uintptr(6)},
			},
			want: `{"level":"info","context":{"uintFieldKey":1,"uint8FieldKey":2,"uint16FieldKey":3,"uint32FieldKey":4,"uint64FieldKey":5,"uintptrFieldKey":6}}`,
		},
		"sets string field": {
			fields: []Field{{Type: StringType, Key: "stringFieldKey", Value: "some string"}},
			want:   `{"level":"info","context":{"stringFieldKey":"some string"}}`,
		},
		"sets time field": {
			fields: []Field{{Type: TimeType, Key: "timeFieldKey", Value: now}},
			want:   fmt.Sprintf(`{"level":"info","context":{"timeFieldKey":"%s"}}`, now.Format(time.RFC3339)),
		},
		"handles unknown field type": {
			fields: []Field{{Type: UnkownType, Key: "unknownFieldKey", Value: struct{ thing int }{thing: 123}}},
			want:   `{"level":"info","context":{"unknownFieldKey":{}}}`,
		},
	}

	for testName, testCase := range testCases {
		tn, tc := testName, testCase

		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			var buffer bytes.Buffer
			adapter := zerologAdapter{logger: zerolog.New(&buffer)}
			adapter.Adapt(tc.level, tc.msg, tc.fields...)
			assert.Equal(t, tc.want+"\n", buffer.String())
		})
	}
}

func TestToZerologLevel(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		level Level
		want  zerolog.Level
	}{
		"debug":            {level: DebugLevel, want: zerolog.DebugLevel},
		"info":             {level: InfoLevel, want: zerolog.InfoLevel},
		"warn":             {level: WarnLevel, want: zerolog.WarnLevel},
		"error":            {level: ErrorLevel, want: zerolog.ErrorLevel},
		"unknown/fallback": {level: Level(123), want: zerolog.DebugLevel},
	}

	for testName, testCase := range testCases {
		tn, tc := testName, testCase

		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, toZerologLevel(tc.level))
		})
	}
}
