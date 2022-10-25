// Package lgrtest exposes an adapter that can be used to assert what logs have been written during a test run.
package lgrtest

import (
	"github.com/stretchr/testify/assert"

	"github.com/nickbryan/collectable/libraries/lgr"
)

// New constructs a new lgr.Logger that has a test adapter for capturing all log entries
// within the returned Entries object. This logger can be a direct replacement for an application
// logger so logs can be asserted when writing automated tests.
func New() (*lgr.Logger, *Entries) {
	e := &Entries{
		entries: []Entry{},
	}

	return lgr.FromAdapter(testAdapter{entries: e}), e
}

// Entry represents a single log entry that was created by the logger.
type Entry struct {
	// Msg is the message that was written in this log entry.
	Msg string
	// Level is the level that this log entry was written at.
	Level lgr.Level
	// Fields are the contextual fields that were added to this log entry.
	Fields map[string]lgr.Field
}

// Entries is a object that allows access to the entries that were logged via the logger.
type Entries struct {
	entries []Entry
}

// All will return all Entry objects that were logged via the logger.
func (e *Entries) All() []Entry {
	return e.entries
}

// Idx is short for entries.All()[idx].
func (e *Entries) Idx(idx uint) Entry {
	return e.entries[idx]
}

type testAdapter struct {
	entries *Entries
}

func (t testAdapter) Adapt(level lgr.Level, message string, fields ...lgr.Field) {
	mappedFields := make(map[string]lgr.Field, len(fields))

	for _, f := range fields {
		mappedFields[f.Key] = f
	}

	t.entries.entries = append(t.entries.entries, Entry{
		Msg:    message,
		Level:  level,
		Fields: mappedFields,
	})
}

// TestingT allows us to abstract away from testing.T so we can properly test our assertion helpers.
type TestingT interface {
	Errorf(format string, args ...interface{})
}

// AssertFullEntry is a test assertion helper that wraps up common assertions for checking a full log entry.
// In order for this assertion to pass:
//
//	entry.Level must match expectedLevel
//	entry.Msg must match expectedMessage
//	len(entry.Fields) must match len(expectedFields)
//	all expectedFields must exist in entry.Fields which is checked by Key
//	each expectedFields Type must match the Type of the matching entry Field
//	each expectedFields Value must match the Value of the matching entry Field
func AssertFullEntry(
	t TestingT, //nolint: varnamelen // t is descriptive of the type.
	entry Entry,
	expectedLevel lgr.Level,
	expectedMessage string,
	expectedFields ...lgr.Field,
) {
	if h, ok := t.(interface { // If we have a testing.T and not a mocked TestingT then tell Go this is a helper func.
		Helper()
	}); ok {
		h.Helper()
	}

	assert.Equal(t, expectedLevel, entry.Level, "level does not match expected")
	assert.Equal(t, expectedMessage, entry.Msg, "message does not match expected")
	assert.Len(t, entry.Fields, len(expectedFields), "length of entry fields does not match length of expected fields")

	for _, expectedField := range expectedFields {
		if field, ok := entry.Fields[expectedField.Key]; ok {
			assert.Equal(
				t,
				expectedField.Type,
				field.Type,
				"field.Type does not match expected for field with key %s",
				expectedField.Key,
			)
			assert.Equal(
				t,
				expectedField.Value,
				field.Value,
				"field.Value does not match expected for field with key %s",
				expectedField.Key,
			)
		} else {
			assert.Fail(t, "field does not exist", "expectedField: %s", expectedField.Key)
		}
	}
}
