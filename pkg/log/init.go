package log

import (
	"github.com/rs/zerolog"
)

var (
	Logger zerolog.Logger
)

func InitLog() {
	Logger = zerolog.New(zerolog.NewConsoleWriter()).With().
		Timestamp().Caller().Logger()
}
