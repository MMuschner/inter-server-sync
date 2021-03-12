package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)



func WriteLog(error error, level zerolog.Level, message string) string{
	const layout = "01-02-2006"
	now := time.Now()
	lf := "uyuni_iss_log_" + now.Format(layout) + ".json"
	logfile, err := os.OpenFile(lf, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if os.IsNotExist(err) {
		f, err := os.Create(lf)
		if err == nil {
		log.Error().Msg("Error: Logfile could not be created.")
		logfile = f
		}
	}
	defer logfile.Close()
	logger := zerolog.New(logfile).With().Timestamp().Logger()

	switch l := level; l {
	case zerolog.InfoLevel:
		logger.Info().Msg(message)
	case zerolog.ErrorLevel:
		logger.Err(error)
	case zerolog.DebugLevel:
		logger.Debug().Msg(message)
	case zerolog.FatalLevel:
		logger.Fatal().Msg(message)
	}

	return lf
}
