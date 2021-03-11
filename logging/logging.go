package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"time"
)

var logfile io.Writer

func setup() string{
	const layout = "01-02-2006"
	now := time.Now()
	lf := "uyuni_iss_log_" + now.Format(layout) + ".json"
	// commented out for testing purposes logfile := "/var/log/rhn/uyuni_iss_log_" + now.Format(layout) + ".json"
	f, err := os.Create(lf)
	if err != nil {
		log.Info().Msg("Error: Logfile could not be created.")
	}
	defer f.Close()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// log := zerolog.New(zerolog.ConsoleWriter{Out: lf, NoColor: false})
	logfile = f
	log := zerolog.New(logfile)
	log.Info().Msg("Test, toast")
	return lf
}


func WriteLog()string {
	logfile := setup()


	return logfile
}
