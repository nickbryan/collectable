package lgr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nickbryan/collectable/libraries/lgr"
	"github.com/nickbryan/collectable/libraries/lgr/lgrtest"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("creates a logger", func(t *testing.T) {
		t.Parallel()

		logger, err := lgr.New()
		assert.NoError(t, err)
		assert.IsType(t, &lgr.Logger{}, logger)
	})
}

func TestLogger(t *testing.T) {
	t.Parallel()

	logger, entries := lgrtest.New()

	logger.Debug("my debug message", lgr.Str("strKey", "some string"), lgr.Integer("intKey", 123))
	logger.Info("my info message", lgr.Str("strKey", "some string"), lgr.Integer("intKey", 123))
	logger.Warn("my warn message", lgr.Str("strKey", "some string"), lgr.Integer("intKey", 123))
	logger.Error("my error message", lgr.Str("strKey", "some string"), lgr.Integer("intKey", 123))

	logs := entries.All()
	assert.Len(t, logs, 4)

	lgrtest.AssertFullEntry(t, logs[0], lgr.DebugLevel, "my debug message", lgr.Str("strKey", "some string"), lgr.Integer("intKey", 123))
	lgrtest.AssertFullEntry(t, logs[1], lgr.InfoLevel, "my info message", lgr.Str("strKey", "some string"), lgr.Integer("intKey", 123))
	lgrtest.AssertFullEntry(t, logs[2], lgr.WarnLevel, "my warn message", lgr.Str("strKey", "some string"), lgr.Integer("intKey", 123))
	lgrtest.AssertFullEntry(t, logs[3], lgr.ErrorLevel, "my error message", lgr.Str("strKey", "some string"), lgr.Integer("intKey", 123))
}

func TestNewNop(t *testing.T) {
	t.Parallel()

	logger := lgr.NewNop()

	// No assertions as nil logger would panic on receiver normally causing the test to fail.
	logger.Debug("my debug message", lgr.Str("strKey", "some string"))
	logger.Info("my info message", lgr.Str("strKey", "some string"))
	logger.Warn("my warn message", lgr.Str("strKey", "some string"))
	logger.Error("my error message", lgr.Str("strKey", "some string"))
}
