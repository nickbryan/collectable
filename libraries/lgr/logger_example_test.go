package lgr_test

import (
	"log"
	"time"

	"github.com/nickbryan/collectable/libraries/lgr"
)

func Example() {
	logger, err := lgr.New(
		lgr.WithOutputPath("stdout"),
		lgr.WithTimestampFactory(func() time.Time { return time.Date(2022, time.March, 5, 0, 0, 0, 0, time.UTC) }),
		lgr.WithMinLevel(lgr.InfoLevel),
	)
	if err != nil {
		log.Fatalf("creating logger: %v", err)
	}

	logger.Debug("my debug message", lgr.Str("myStringKey", "my string value"))
	logger.Info("my info message", lgr.Integer("myIntKey", 123456))
	logger.Warn("my warn message", lgr.Float("myFloatKey", 12.3))
	logger.Error("my error message", lgr.Time("myTimeKey", time.Date(2021, time.February, 1, 0, 0, 0, 0, time.UTC)))

	// Output:
	// {"level":"info","context":{"myIntKey":123456},"timestamp":"2022-03-05T00:00:00Z","message":"my info message"}
	// {"level":"warn","context":{"myFloatKey":12.3},"timestamp":"2022-03-05T00:00:00Z","message":"my warn message"}
	// {"level":"error","context":{"myTimeKey":"2021-02-01T00:00:00Z"},"timestamp":"2022-03-05T00:00:00Z","message":"my error message"}
}
