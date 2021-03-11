package logging

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)



func WriteLog(error error) string{
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
	logger := zerolog.New(logfile)
	logger.Info().Msg("Test, tea")
	logger.Info().Err(error)

	return lf
}



/*
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

*/
