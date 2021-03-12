package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"time"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/inter-server-sync/cli"
	"github.com/uyuni-project/inter-server-sync/dumper"
	"github.com/uyuni-project/inter-server-sync/schemareader"
)

// func init() { log.SetFlags(log.Lshortfile | log.LstdFlags) }

var Logfile string

func main() {
	const layout = "01-02-2006"
	now := time.Now()
	Logfile := "uyuni_iss_log_" + now.Format(layout) + ".json"
	lf, err := os.OpenFile(Logfile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if os.IsNotExist(err) {
		f, err := os.Create(Logfile)
		if err != nil {
			log.Error().Msg("Unable to create logfile")
		}
		lf = f
	}
	logger := zerolog.New(lf).With().Timestamp().Logger().Output(lf)
	parsedArgs, err := cli.CliArgs(os.Args)
	if err != nil {
		logger.Fatal().Err(err).Msg("")
		os.Exit(1)
	}

	if parsedArgs.Cpuprofile != "" {
		f, err := os.Create(parsedArgs.Cpuprofile)
		if err != nil {
			logger.Fatal().Err(err).Msg("Could not create CPU profile")
		}
		logger.Info().Msg("CPU profile parsed")
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			logger.Fatal().Err(err).Msg("Could not start CPU profile")
		}
		defer pprof.StopCPUProfile()
	}

	db := schemareader.GetDBconnection(parsedArgs.Config)
	defer db.Close()

	if parsedArgs.Dot {
		tables := schemareader.ReadTablesSchema(db, dumper.SoftwareChannelTableNames())
		schemareader.DumpToGraphviz(tables)
		return
	}
	if len(parsedArgs.ChannleLabels) > 0 {
		channelLabels := parsedArgs.ChannleLabels
		outputFolder := parsedArgs.Path
		tableData := dumper.DumpChannelData(db, channelLabels, outputFolder)

		if parsedArgs.Debug {
			for index, channelTableData := range tableData {
				fmt.Printf("###Processing channe%d...", index)
				for path := range channelTableData.Paths {
					println(path)
				}
				count := 0
				for _, value := range channelTableData.TableData {
					fmt.Printf("%s number inserts: %d \n\t %s keys: %s\n", value.TableName, len(value.Keys),
						value.TableName, value.Keys)
					count = count + len(value.Keys)
				}
				fmt.Printf("IDS############%d\n\n", count)
			}

		}
	}
	if parsedArgs.Memprofile != "" {
		f, err := os.Create(parsedArgs.Memprofile)
		if err != nil {
			logger.Fatal().Err(err).Msg("Could not create memory profile")
		}
		defer f.Close() // error handling omitted for example
		//runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			logger.Fatal().Err(err).Msg("Could not write memory profile")
		}
	}

}
