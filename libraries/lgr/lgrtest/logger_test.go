package lgrtest_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nickbryan/collectable/libraries/lgr"
	"github.com/nickbryan/collectable/libraries/lgr/lgrtest"
)

// bufferT implements TestingT. Its implementation of Errorf writes the output that would be produced by
// testing.T.Errorf to an internal bytes.Buffer.
type bufferT struct {
	buf bytes.Buffer
}

func (t *bufferT) Errorf(format string, args ...interface{}) {
	t.buf.WriteString(fmt.Sprintf(format, args...))
}

func TestAssertFullEntry(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		entry          lgrtest.Entry
		expectedLevel  lgr.Level
		expectedMsg    string
		expectedFields []lgr.Field
		wantMsg        string
	}{
		"level does not match": {
			entry:         lgrtest.Entry{Level: lgr.DebugLevel},
			expectedLevel: lgr.InfoLevel,
			wantMsg:       "level does not match expected",
		},
		"message does not match": {
			entry:         lgrtest.Entry{Level: lgr.DebugLevel, Msg: "some log message"},
			expectedLevel: lgr.DebugLevel,
			expectedMsg:   "some other log message",
			wantMsg:       "message does not match expected",
		},
		"fields len do not match": {
			entry:          lgrtest.Entry{Level: lgr.DebugLevel, Msg: "some log message"},
			expectedLevel:  lgr.DebugLevel,
			expectedMsg:    "some log message",
			expectedFields: []lgr.Field{lgr.Str("someKey", "someVal")},
			wantMsg:        "length of entry fields does not match length of expected fields",
		},
		"field not found": {
			entry: lgrtest.Entry{
				Level:  lgr.DebugLevel,
				Msg:    "some log message",
				Fields: map[string]lgr.Field{"someKey": lgr.Str("someKey", "someVal")},
			},
			expectedLevel:  lgr.DebugLevel,
			expectedMsg:    "some log message",
			expectedFields: []lgr.Field{lgr.Str("someOtherKey", "someVal")},
			wantMsg:        "expectedField: someOtherKey",
		},
		"field type does not match": {
			entry: lgrtest.Entry{
				Level:  lgr.DebugLevel,
				Msg:    "some log message",
				Fields: map[string]lgr.Field{"someKey": lgr.Str("someKey", "someVal")},
			},
			expectedLevel:  lgr.DebugLevel,
			expectedMsg:    "some log message",
			expectedFields: []lgr.Field{lgr.Bool("someKey", false)},
			wantMsg:        "field.Value does not match expected for field with key someKey",
		},
		"field value does not match": {
			entry: lgrtest.Entry{
				Level:  lgr.DebugLevel,
				Msg:    "some log message",
				Fields: map[string]lgr.Field{"someKey": lgr.Str("someKey", "someVal")},
			},
			expectedLevel:  lgr.DebugLevel,
			expectedMsg:    "some log message",
			expectedFields: []lgr.Field{lgr.Str("someKey", "someOtherVal")},
			wantMsg:        "field.Value does not match expected for field with key someKey",
		},
	}

	for testName, testCase := range testCases {
		tn, tc := testName, testCase

		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			mockT := new(bufferT)
			lgrtest.AssertFullEntry(mockT, tc.entry, tc.expectedLevel, tc.expectedMsg, tc.expectedFields...)
			assert.Contains(t, mockT.buf.String(), tc.wantMsg)
		})
	}
}
