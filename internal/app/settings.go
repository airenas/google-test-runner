package app

import (
	"errors"
	"flag"
	"os"
)

//Settings for the app
type Settings struct {
	workingDir     string
	showOnlyFailed bool
	workersCount   int
	reader         reader
	formatter      formatter
}

//InitSettings initializes the app settings
func InitSettings() (*Settings, error) {
	filePtr := flag.String("l", "", "File of google executable list, if empty stdin. Sample: -l files.in")
	filterPtr := flag.String("f", "", "Filter to search recursivelly for files in current dir. Sample: -f ./**/*_test")
	workingDirPtr := flag.String("d", "./", "Working dir: Sample: -d ./")
	wCount := flag.Int("j", 4, "Workers count to run in parallel. Allowed value: [1, 99]. Sample -j 4")
	onlyFailed := flag.Bool("s", false, "Show only failed test cases")
	showGTestOutput := flag.Bool("o", false, "Show original GTest output")
	flag.Parse()
	result := Settings{}
	result.workingDir = *workingDirPtr
	var err error
	result.reader, err = initReader(*filePtr, *filterPtr, *workingDirPtr)
	if err != nil {
		return nil, err
	}
	result.formatter = newFormatterWriter(os.Stdout, *onlyFailed, *showGTestOutput)
	result.showOnlyFailed = *onlyFailed
	result.workersCount = *wCount

	if result.workersCount < 1 || result.workersCount > 100 {
		return nil, errors.New("Workers count must be in [1, 99]. Run: googleTestRunner -h")
	}

	return &result, nil
}

func initReader(fileList string, filter string, wDir string) (reader, error) {
	if fileList != "" && filter != "" {
		return nil, errors.New("Only l or f parameter allowed, not both")
	}
	if filter != "" {
		return newDirReader(wDir, filter), nil
	}
	if fileList != "" {
		return newFileReader(fileList), nil
	}
	return newStdInReader(fileList), nil

}
