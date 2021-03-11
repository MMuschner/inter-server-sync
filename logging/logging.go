package logging

import (
	"fmt"
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
	 logger := zerolog.New(zerolog.ConsoleWriter{Out: f, NoColor: false})
	 // logger := zerolog.New(f)
	 logger.Info().Msg("Test, toast")
	return lf
}


func WriteLog(error error)string {
	logfile := setup()
	lf, err := os.Open(logfile)
	logger := zerolog.New(lf)
	if err != nil {
		fmt.Println("error found")
		logger.Info().Msg("Error handling logfile.")
	}
	logger.Info().Msg("Successful entry")
	fmt.Println("Error is: ", error)
	fmt.Println("Logfile is: ", logfile)
	defer lf.Close()
	return logfile
}
