package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strings"

	"github.com/harry1453/go-common-file-dialog/cfd"
	"github.com/harry1453/go-common-file-dialog/cfdutil"

	"github.com/genshinsim/gcsim/pkg/core"
	"github.com/genshinsim/gcsim/pkg/parse"
	"github.com/genshinsim/gcsim/pkg/result"
	"github.com/genshinsim/gcsim/simulator"
	"github.com/genshinsim/gcsim/substatoptimizer"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	dataframe "github.com/rocketlaunchr/dataframe-go"
	exports "github.com/rocketlaunchr/dataframe-go/exports"
)

type dpsPerChar struct {
	Char string
	Dps  float64
}

func parseConfig(simopt simulator.Options) (float64, float64, []string, []float64) {
	// Parse config
	zapcfg := zap.NewDevelopmentConfig()
	zapcfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	zapcfg.EncoderConfig.CallerKey = ""
	zapcfg.EncoderConfig.StacktraceKey = ""
	zapcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// if verbose {
	// 	zapcfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	// }
	logger, _ := zapcfg.Build()
	defer logger.Sync()
	sugarLog := logger.Sugar()
	cfg, err := simulator.ReadConfig(simopt.ConfigPath)
	if err != nil {
		sugarLog.Error(err)
		os.Exit(1)
	}

	//get the characters names to store them later in array
	// var reGetCharNames = regexp.MustCompile(`(?m)^([a-z]+)\s+char\b[^;]*;`)
	// for _, match := range reGetCharNames.FindAllStringSubmatch(cfg, -1) {
	// 	char := string(match[1])
	// 	fmt.Printf("%q\n", char)

	// }
	// if err != nil {
	// 	log.Println(err)
	// 	os.Exit(1)
	// }
	parser := parse.New("single", string(cfg))
	simcfg, err := parser.Parse()
	result := runSimWithConfig(cfg, simcfg, simopt)
	//fmt.Printf("DPS: %v     STD DEV: %v \n", result.DPS.Mean, result.DPS.SD)
	total := make([]float64, len(result.CharNames), len(result.CharNames))

	for i, t := range result.DamageByChar {

		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}

		for _, k := range keys {
			v := t[k]

			total[i] += v.Mean
		}

	}

	return result.DPS.Mean, result.DPS.SD, result.CharNames, total
}

// Just runs the sim with specified settings
func runSimWithConfig(cfg string, simcfg core.SimulationConfig, simopt simulator.Options) result.Summary {
	result, err := simulator.RunWithConfig(cfg, simcfg, simopt)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return result
}

func selectFiles() []string {
	results, err := cfdutil.ShowOpenMultipleFilesDialog(cfd.DialogConfig{
		Title: "Select the configs you want to run",
		Role:  "OpenFilesExample",
		FileFilters: []cfd.FileFilter{
			{
				DisplayName: "Text Files (*.txt)",
				Pattern:     "*.txt",
			},
		},
		SelectedFileFilterIndex: 2,
		FileName:                "file.txt",
		DefaultExtension:        "txt",
	})
	if err == cfd.ErrorCancelled {
		log.Fatal("Dialog was cancelled by the user.")
	} else if err != nil {
		log.Fatal(err)
	}

	return results

}

func (d dpsPerChar) String() string { //Format the struct to a pretty form
	return fmt.Sprintf("%s:\n%.2f", d.Char, d.Dps)
}

func OptnRun(optimize bool) {
	filepaths := selectFiles()
	log.Printf("Chosen file(s): %s\n", filepaths)

	s1 := dataframe.NewSeriesString("File name", nil)
	s2 := dataframe.NewSeriesFloat64("Total DPS", nil)
	s3 := dataframe.NewSeriesFloat64("Std Dev", nil)
	sg1 := dataframe.NewSeriesGeneric("Char 1 DPS", dpsPerChar{}, nil)
	sg2 := dataframe.NewSeriesGeneric("Char 2 DPS", dpsPerChar{}, nil)
	sg3 := dataframe.NewSeriesGeneric("Char 3 DPS", dpsPerChar{}, nil)
	sg4 := dataframe.NewSeriesGeneric("Char 4 DPS", dpsPerChar{}, nil)
	var ctx = context.Background()
	var excelExportOpt exports.ExcelExportOptions

	df := dataframe.NewDataFrame(s1, s2, s3, sg1, sg2, sg3, sg4)

	for _, filepath := range filepaths {
		var configpath simulator.Options
		configpath.ConfigPath = filepath
		configpath.GZIPResult = true //saves .gz
		configpath.ResultSaveToPath = filepath
		filenamepath := strings.Split(configpath.ResultSaveToPath, "\\")
		filename := filenamepath[len(filenamepath)-1]

		if optimize {
			log.Printf("Optimizing: %s ...", filepath)
			substatoptimizer.RunSubstatOptim(configpath, false, "")
		}

		log.Printf("Simulating: %s ...", filepath)
		configpath.ResultSaveToPath = strings.Replace(filepath, ".txt", "", -1)
		dpsTotal, sd, chars, dpsChar := parseConfig(configpath) //Calls the sim to run
		fmt.Printf("Done!\n")

		df.Append(nil, filename, math.Round(dpsTotal*100)/100, math.Round(sd*100)/100, dpsPerChar{chars[0], dpsChar[0]},
			dpsPerChar{chars[1], dpsChar[1]},
			dpsPerChar{chars[2], dpsChar[2]},
			dpsPerChar{chars[3], dpsChar[3]})
	}
	sortDps := []dataframe.SortKey{
		{Key: "Total DPS", Desc: true},
	}
	df.Sort(ctx, sortDps)
	fmt.Print(df.Table())

	writer, err := os.Create("Results.xlsx")
	if err != nil {
		fmt.Print(err)
	}

	errExport := exports.ExportToExcel(ctx, writer, df, excelExportOpt)
	if errExport != nil {
		fmt.Print(errExport)
	} else {
		fmt.Printf("Saved to results.xlsx!\n")

	}
	fmt.Println("\nPress the Enter Key to close program!")
	fmt.Scanln()
}
