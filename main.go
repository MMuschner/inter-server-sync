package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/uyuni-project/inter-server-sync/cli"
	"github.com/uyuni-project/inter-server-sync/dumper"
	"github.com/uyuni-project/inter-server-sync/schemareader"
	"github.com/uyuni-project/inter-server-sync/logging"
	"os"
	"runtime/pprof"
)
/*
// func init() { log.SetFlags(log.Lshortfile | log.LstdFlags) }
func init() {
	const layout = "01-02-2006"
	now := time.Now()
	logfile := "uyuni_iss_log_" + now.Format(layout) + ".json"
	// commented out for testing purposes logfile := "/var/log/rhn/uyuni_iss_log_" + now.Format(layout) + ".json"
	lf, err := os.Create(logfile)
	if err != nil {
		log.Info().Msg("Error: Logfile could not be created.")
	}
	defer lf.Close()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// log := zerolog.New(zerolog.ConsoleWriter{Out: lf, NoColor: false})
	log := zerolog.New(lf)
	log.Info().Msg("Test")

}

 */


func main() {
	parsedArgs, err := cli.CliArgs(os.Args)
	if err != nil {
		logging.WriteLog()
		os.Exit(1)
	}

	if parsedArgs.Cpuprofile != "" {
		f, err := os.Create(parsedArgs.Cpuprofile)
		if err != nil {
			log.Fatal().Err(err).Msgf("Could not create CPU profile")
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			// log.Fatal("could not start CPU profile: ", err)

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
			// log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		//runtime.GC()    // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			// log.Fatal("could not write memory profile: ", err)
		}
	}

}
