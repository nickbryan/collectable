package lgr

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

type zerologAdapter struct {
	logger zerolog.Logger
}

func bindZerologAdapter(l *Logger) {
	var output io.Writer

	switch l.outputPath {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	default:
		os.OpenFile(l.outputPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	}

	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimestampFunc = l.timestampFactory
	logger := zerolog.New(output).Level(toZerologLevel(l.minLevel)).With().Timestamp().Logger()

	l.adapter = zerologAdapter{logger: logger}
}

func (z zerologAdapter) Adapt(level Level, message string, fields ...Field) {
	var event *zerolog.Event

	switch level {
	case DebugLevel:
		event = z.logger.Debug()
	case InfoLevel:
		event = z.logger.Info()
	case WarnLevel:
		event = z.logger.Warn()
	case ErrorLevel:
		event = z.logger.Error()
	default:
		panic(fmt.Sprintf("log level unexpected: %v", level))
	}

	event.Dict("context", fieldsToContext(fields))

	event.Msg(message)
}
func fieldsToContext(fields []Field) *zerolog.Event { //nolint: cyclop,funlen // Easier to read whole type switch.
	event := zerolog.Dict()

	for _, field := range fields {
		switch field.Type {
		case BoolType:
			event.Bool(field.Key, field.Value.(bool)) //nolint: forcetypeassert // We know the type.
		case ByteStringType:
			event.Bytes(field.Key, field.Value.([]byte)) //nolint: forcetypeassert // We know the type.
		case DurationType:
			event.Dur(field.Key, field.Value.(time.Duration)) //nolint: forcetypeassert // We know the type.
		case ErrorType:
			event.AnErr(field.Key, field.Value.(error)) //nolint: forcetypeassert // We know the type.
		case Float32Type:
			event.Float32(field.Key, field.Value.(float32)) //nolint: forcetypeassert // We know the type.
		case Float64Type:
			event.Float64(field.Key, field.Value.(float64)) //nolint: forcetypeassert // We know the type.
		case IntType:
			event.Int(field.Key, field.Value.(int)) //nolint: forcetypeassert // We know the type.
		case Int8Type:
			event.Int8(field.Key, field.Value.(int8)) //nolint: forcetypeassert // We know the type.
		case Int16Type:
			event.Int16(field.Key, field.Value.(int16)) //nolint: forcetypeassert // We know the type.
		case Int32Type:
			event.Int32(field.Key, field.Value.(int32)) //nolint: forcetypeassert // We know the type.
		case Int64Type:
			event.Int64(field.Key, field.Value.(int64)) //nolint: forcetypeassert // We know the type.
		case UintType:
			event.Uint(field.Key, field.Value.(uint)) //nolint: forcetypeassert // We know the type.
		case Uint8Type:
			event.Uint8(field.Key, field.Value.(uint8)) //nolint: forcetypeassert // We know the type.
		case Uint16Type:
			event.Uint16(field.Key, field.Value.(uint16)) //nolint: forcetypeassert // We know the type.
		case Uint32Type:
			event.Uint32(field.Key, field.Value.(uint32)) //nolint: forcetypeassert // We know the type.
		case Uint64Type:
			event.Uint64(field.Key, field.Value.(uint64)) //nolint: forcetypeassert // We know the type.
		case UintptrType:
			event.Uint64(field.Key, uint64(field.Value.(uintptr))) //nolint: forcetypeassert // We know the type.
		case StringType:
			event.Str(field.Key, field.Value.(string)) //nolint: forcetypeassert // We know the type.
		case TimeType:
			event.Time(field.Key, field.Value.(time.Time)) //nolint: forcetypeassert // We know the type.
		case UnkownType:
			fallthrough
		default:
			event.Interface(field.Key, field.Value)
		}
	}

	return event
}

func toZerologLevel(l Level) zerolog.Level {
	switch l {
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	default:
		return zerolog.DebugLevel // Fallback to DebugLevel for full log output.
	}
}
