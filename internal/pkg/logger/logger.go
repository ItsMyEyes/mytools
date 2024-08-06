package logger

import (
	"os"
	"strconv"
	"sync"

	"github.com/gobuffalo/envy"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

const (
	// LogLevelEnv is an environment variable name for LOG_LEVEL
	LogLevelEnv = "LOG_LEVEL"
	// EnableJSONLogsEnv is an environment variable name for ENABLE_JSON_LOGS
	EnableJSONLogsEnv = "ENABLE_JSON_LOGS"
)

var once sync.Once

var log zerolog.Logger

func Get() zerolog.Logger {
	once.Do(func() {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		logLevel, err := strconv.Atoi(envy.Get(LogLevelEnv, "0"))
		if err != nil {
			logLevel = int(zerolog.InfoLevel)
		}

		log = zerolog.New(os.Stdout).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Caller().
			Logger()

		if envy.Get(EnableJSONLogsEnv, "false") == "false" {
			log = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
		}

		zerolog.DefaultContextLogger = &log
		zlog.Logger = log

	})
	return log
}
